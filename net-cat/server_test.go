package main

import (
	"bufio"
	"io"
	"net"
	"os"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
)

func testServer() *server {
	return &server{
		g:  make(map[string]*group),
		bq: make(chan message, 16),
	}
}

func testPipeClient(t *testing.T, s *server, name string) (*client, net.Conn, *bufio.Reader) {
	t.Helper()

	serverConn, peerConn := net.Pipe()
	t.Cleanup(func() {
		serverConn.Close()
		peerConn.Close()
	})
	return &client{
		s:      s,
		fd:     serverConn,
		reader: bufio.NewReader(serverConn),
		name:   name,
	}, peerConn, bufio.NewReader(peerConn)
}

func writeTestPeer(t *testing.T, peer net.Conn, text string) {
	t.Helper()
	if err := peer.SetWriteDeadline(time.Now().Add(2 * time.Second)); err != nil {
		t.Fatal(err)
	}
	if _, err := io.WriteString(peer, text); err != nil {
		t.Fatal(err)
	}
}

func readTestLine(t *testing.T, peer net.Conn, reader *bufio.Reader) string {
	t.Helper()
	if err := peer.SetReadDeadline(time.Now().Add(2 * time.Second)); err != nil {
		t.Fatal(err)
	}
	line, err := reader.ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}
	return line
}

func readTestDelimiter(t *testing.T, peer net.Conn, reader *bufio.Reader, delim byte) string {
	t.Helper()
	if err := peer.SetReadDeadline(time.Now().Add(2 * time.Second)); err != nil {
		t.Fatal(err)
	}
	text, err := reader.ReadString(delim)
	if err != nil {
		t.Fatal(err)
	}
	return text
}

func TestConnectionLimitReservations(t *testing.T) {
	s := server{}

	for i := 0; i < maxClients; i++ {
		if !s.reserveConnection() {
			t.Fatalf("reservation %d was rejected before the limit", i+1)
		}
	}

	if s.reserveConnection() {
		t.Fatal("reservation exceeded the connection limit")
	}

	s.releaseConnection()
	if !s.reserveConnection() {
		t.Fatal("reservation was rejected after a connection slot was released")
	}

	if s.reserveConnection() {
		t.Fatal("reservation exceeded the connection limit after slot reuse")
	}
}

func TestConcurrentConnectionReservationsRespectLimit(t *testing.T) {
	s := server{}
	const attempts = maxClients * 10
	results := make(chan bool, attempts)

	var wg sync.WaitGroup
	for i := 0; i < attempts; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			results <- s.reserveConnection()
		}()
	}
	wg.Wait()
	close(results)

	reserved := 0
	for ok := range results {
		if ok {
			reserved++
		}
	}
	if reserved != maxClients {
		t.Fatalf("reserved %d connections, want %d", reserved, maxClients)
	}
}

func TestParseArgs(t *testing.T) {
	originalArgs := os.Args
	t.Cleanup(func() {
		os.Args = originalArgs
	})

	tests := []struct {
		name    string
		args    []string
		want    string
		wantErr string
	}{
		{name: "default port", args: []string{"TCPChat"}, want: ":8989"},
		{name: "specified port", args: []string{"TCPChat", "2525"}, want: ":2525"},
		{
			name:    "extra argument",
			args:    []string{"TCPChat", "2525", "localhost"},
			wantErr: "[USAGE]: ./TCPChat $port",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.args
			got, err := parseArgs()
			if tt.wantErr != "" {
				if err == nil || err.Error() != tt.wantErr {
					t.Fatalf("parseArgs() error = %v, want %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if got.addr != tt.want {
				t.Fatalf("parseArgs() addr = %q, want %q", got.addr, tt.want)
			}
		})
	}
}

func TestAcceptConnection(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	s := &server{sock: listener}
	dialDone := make(chan error, 1)
	go func() {
		conn, err := net.Dial("tcp", listener.Addr().String())
		if err == nil {
			defer conn.Close()
		}
		dialDone <- err
	}()

	c, err := acceptConnection(s)
	if err != nil {
		t.Fatal(err)
	}
	defer c.fd.Close()

	if c.s != s {
		t.Fatal("accepted client does not reference its server")
	}
	if c.reader == nil {
		t.Fatal("accepted client has no buffered reader")
	}
	if c.fd.RemoteAddr() == nil {
		t.Fatal("accepted connection has no remote address")
	}
	if err := <-dialDone; err != nil {
		t.Fatal(err)
	}
}

func TestInitClientBannerNamePromptAndMembership(t *testing.T) {
	s := testServer()
	s.g[""] = &group{name: "", history: "earlier message\n"}
	c, peer, reader := testPipeClient(t, s, "")
	done := make(chan error, 1)
	go func() {
		done <- initClient(c)
	}()

	banner := readTestDelimiter(t, peer, reader, ':')
	if !strings.Contains(banner, "Welcome to TCP-Chat!") ||
		!strings.Contains(banner, "_nnnn_") ||
		!strings.HasSuffix(banner, "[ENTER YOUR NAME]:") {
		t.Fatalf("unexpected welcome banner: %q", banner)
	}
	if space, err := reader.ReadByte(); err != nil || space != ' ' {
		t.Fatalf("welcome prompt trailing byte = %q, error = %v", space, err)
	}

	writeTestPeer(t, peer, " \t \n")
	prompt := readTestDelimiter(t, peer, reader, ':')
	if !strings.HasSuffix(prompt, "[ENTER YOUR NAME]:") {
		t.Fatalf("empty name prompt = %q, want another name prompt", prompt)
	}

	readTestDelimiter(t, peer, reader, 'm')
	writeTestPeer(t, peer, "Alice\n")
	if got := readTestLine(t, peer, reader); got != "earlier message\n" {
		t.Fatalf("new client's history = %q, want earlier message", got)
	}
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	if c.name != "Alice" {
		t.Fatalf("client name = %q, want Alice", c.name)
	}
	if c.g == nil || len(c.g.c) != 1 {
		t.Fatal("named client was not added to a group")
	}
}

func TestInitClientReturnsErrorWhenConnectionCloses(t *testing.T) {
	s := testServer()
	c, peer, _ := testPipeClient(t, s, "")
	peer.Close()

	if err := initClient(c); err == nil {
		t.Fatal("initClient succeeded after the connection closed")
	}
}

func TestClientGroupAddOrdersHistoryBeforeLiveMessages(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	defer serverConn.Close()
	defer clientConn.Close()
	if err := clientConn.SetReadDeadline(time.Now().Add(2 * time.Second)); err != nil {
		t.Fatal(err)
	}

	senderConn, senderPeer := net.Pipe()
	defer senderConn.Close()
	defer senderPeer.Close()
	go io.Copy(io.Discard, senderPeer)

	s := &server{g: make(map[string]*group)}
	grp := &group{name: "general", history: "previous message\n"}
	sender := &client{s: s, g: grp, fd: senderConn}
	grp.c = []*client{sender}
	s.g[grp.name] = grp

	joining := &client{s: s, fd: serverConn}
	addDone := make(chan error, 1)
	go func() {
		addDone <- clientGroupAdd(joining, grp.name)
	}()

	deadline := time.Now().Add(time.Second)
	for {
		s.mu.Lock()
		joined := joining.g != nil
		s.mu.Unlock()
		if joined {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("client was not added to the group")
		}
		runtime.Gosched()
	}

	broadcastDone := make(chan struct{})
	go func() {
		broadcastSystemMessage(sender, "live during history")
		close(broadcastDone)
	}()

	reader := bufio.NewReader(clientConn)
	first, err := reader.ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}
	if first != "previous message\n" {
		t.Fatalf("first message = %q, want history before live messages", first)
	}

	second, err := reader.ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(second, "live during history") {
		t.Fatalf("second message = %q, want live message after history", second)
	}

	if err := <-addDone; err != nil {
		t.Fatal(err)
	}
	<-broadcastDone
}

func TestClientGroupAddRemovesClientWhenHistoryWriteFails(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	clientConn.Close()
	defer serverConn.Close()

	s := &server{g: make(map[string]*group)}
	grp := &group{name: "general", history: "previous message\n"}
	s.g[grp.name] = grp
	c := &client{s: s, fd: serverConn}

	if err := clientGroupAdd(c, grp.name); err == nil {
		t.Fatal("clientGroupAdd succeeded after the history write failed")
	}
	if c.g != nil {
		t.Fatal("client remained associated with a group after the history write failed")
	}
	if len(grp.c) != 0 {
		t.Fatalf("group has %d clients after failed history write, want 0", len(grp.c))
	}
	if _, ok := s.g[grp.name]; ok {
		t.Fatal("empty group remained registered after failed history write")
	}
}

func TestClientGroupAddAndDelNotifyOtherMembers(t *testing.T) {
	s := testServer()
	existing, existingPeer, existingReader := testPipeClient(t, s, "Existing")
	joining, joiningPeer, joiningReader := testPipeClient(t, s, "Joining")

	grp := &group{name: "general", history: "earlier message\n"}
	grp.c = []*client{existing}
	existing.g = grp
	s.g[grp.name] = grp

	addDone := make(chan error, 1)
	go func() {
		addDone <- clientGroupAdd(joining, grp.name)
	}()

	if got := readTestLine(t, joiningPeer, joiningReader); got != "earlier message\n" {
		t.Fatalf("joining client history = %q, want previous message", got)
	}
	joinNotice := readTestLine(t, existingPeer, existingReader)
	if !strings.Contains(joinNotice, "Joining") || !strings.Contains(joinNotice, "joined") {
		t.Fatalf("join notice = %q", joinNotice)
	}
	if err := <-addDone; err != nil {
		t.Fatal(err)
	}

	delDone := make(chan error, 1)
	go func() {
		delDone <- clientGroupDel(joining)
	}()
	leaveNotice := readTestLine(t, existingPeer, existingReader)
	if !strings.Contains(leaveNotice, "Joining") || !strings.Contains(leaveNotice, "left") {
		t.Fatalf("leave notice = %q", leaveNotice)
	}
	if err := <-delDone; err != nil {
		t.Fatal(err)
	}

	if joining.g != nil {
		t.Fatal("departed client is still associated with the group")
	}
	if len(grp.c) != 1 || grp.c[0] != existing {
		t.Fatal("remaining client was removed when another client left")
	}
	if s.g[grp.name] != grp {
		t.Fatal("group was removed while a client remained")
	}
}

func TestClientGroupAddRejectsFullGroup(t *testing.T) {
	s := testServer()
	grp := &group{name: "general"}
	s.g[grp.name] = grp
	for i := 0; i < maxClients; i++ {
		grp.c = append(grp.c, &client{s: s, g: grp})
	}

	c, peer, reader := testPipeClient(t, s, "Overflow")
	addDone := make(chan error, 1)
	go func() {
		addDone <- clientGroupAdd(c, grp.name)
	}()

	if got := readTestLine(t, peer, reader); got != "Chat is full. Try again later.\n" {
		t.Fatalf("full-group response = %q", got)
	}
	if err := <-addDone; err == nil {
		t.Fatal("clientGroupAdd succeeded for a full group")
	}
	if c.g != nil {
		t.Fatal("rejected client was added to the group")
	}
	if len(grp.c) != maxClients {
		t.Fatalf("group has %d clients after rejection, want %d", len(grp.c), maxClients)
	}
}

func TestClientLoopIgnoresEmptyMessagesAndHandlesRename(t *testing.T) {
	s := testServer()
	c, peer, reader := testPipeClient(t, s, "Alice")
	other, otherPeer, otherReader := testPipeClient(t, s, "Bob")

	grp := &group{name: "general", c: []*client{c, other}}
	c.g = grp
	other.g = grp
	s.g[grp.name] = grp

	done := make(chan error, 1)
	go func() {
		done <- clientLoop(c)
	}()

	writeTestPeer(t, peer, "\n \t \r\n")
	if len(s.bq) != 0 {
		t.Fatal("empty client messages were queued")
	}

	writeTestPeer(t, peer, "/name\n")
	if got := readTestLine(t, peer, reader); !strings.Contains(got, "Usage: /name <new_name>\n") {
		t.Fatalf("name command usage = %q", got)
	}

	writeTestPeer(t, peer, "/name AliceSmith\n")
	renameNotice := readTestLine(t, otherPeer, otherReader)
	if !strings.Contains(renameNotice, "changed their name to") {
		t.Fatalf("rename notice = %q", renameNotice)
	}

	writeTestPeer(t, peer, "hello everyone\n")
	select {
	case msg := <-s.bq:
		if msg.origin != c || msg.origin.name != "AliceSmith" || !strings.HasSuffix(msg.text, ": hello everyone\n") {
			t.Fatalf("queued message = %#v", msg)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("client message was not queued")
	}

	peer.Close()
	if err := <-done; err == nil {
		t.Fatal("clientLoop returned nil after the client disconnected")
	}
}

func TestBroadcastLoopFormatsAndSendsChatMessages(t *testing.T) {
	s := testServer()
	sender := &client{s: s, name: "Alice"}
	first, firstPeer, firstReader := testPipeClient(t, s, "Bob")
	second, secondPeer, secondReader := testPipeClient(t, s, "Charlie")
	grp := &group{name: "general", c: []*client{sender, first, second}}
	sender.g = grp
	first.g = grp
	second.g = grp
	s.g[grp.name] = grp

	go broadcastLoop(s)
	processMessage(message{origin: sender, text: "hello\r\n"})

	firstMessage := readTestLine(t, firstPeer, firstReader)
	secondMessage := readTestLine(t, secondPeer, secondReader)
	if firstMessage != secondMessage {
		t.Fatalf("clients received different messages: %q and %q", firstMessage, secondMessage)
	}

	start := strings.Index(firstMessage, "[20")
	end := strings.Index(firstMessage, "]")
	if start < 0 || end < 0 {
		t.Fatalf("message has no timestamp/name separator: %q", firstMessage)
	}
	if _, err := time.Parse("2006-01-02 15:04:05", firstMessage[start + 1:end]); err != nil {
		t.Fatalf("message timestamp is invalid: %q: %v", firstMessage, err)
	}

	s.mu.Lock()
	history := grp.history
	s.mu.Unlock()
	if history != firstMessage {
		t.Fatalf("history = %q, want sent message %q", history, firstMessage)
	}
}

func TestHandleDisconnectBeforeJoiningGroupIsSafe(t *testing.T) {
	s := testServer()
	c := &client{s: s}
	if err := handleDisconnect(s, c); err != nil {
		t.Fatal(err)
	}
}

func TestHandleClientDisconnectKeepsOtherClientsConnected(t *testing.T) {
	s := testServer()
	existing, existingPeer, existingReader := testPipeClient(t, s, "Existing")
	grp := &group{name: "", history: "earlier message\n", c: []*client{existing}}
	existing.g = grp
	s.g[grp.name] = grp

	if !s.reserveConnection() {
		t.Fatal("could not reserve a client connection")
	}
	s.wg.Add(1)
	c, peer, reader := testPipeClient(t, s, "")
	finished := make(chan struct{})
	go func() {
		handleClient(c)
		close(finished)
	}()

	banner := readTestDelimiter(t, peer, reader, ':')
	if !strings.Contains(banner, "Welcome to TCP-Chat!") {
		t.Fatalf("unexpected welcome banner: %q", banner)
	}
	if _, err := reader.ReadByte(); err != nil {
		t.Fatal(err)
	}
	readTestDelimiter(t, peer, reader, 'm')
	writeTestPeer(t, peer, "Joining\n")
	if got := readTestLine(t, peer, reader); got != "earlier message\n" {
		t.Fatalf("new client's history = %q, want earlier message", got)
	}
	joinNotice := readTestLine(t, existingPeer, existingReader)
	if !strings.Contains(joinNotice, "Joining") || !strings.Contains(joinNotice, "joined") {
		t.Fatalf("join notice = %q", joinNotice)
	}

	peer.Close()
	leaveNotice := readTestLine(t, existingPeer, existingReader)
	if !strings.Contains(leaveNotice, "Joining") || !strings.Contains(leaveNotice, "left") {
		t.Fatalf("leave notice = %q", leaveNotice)
	}

	select {
	case <-finished:
	case <-time.After(2 * time.Second):
		t.Fatal("client handler did not finish after disconnect")
	}
	s.mu.Lock()
	activeConnections := s.activeConnections
	memberCount := len(grp.c)
	s.mu.Unlock()
	if activeConnections != 0 {
		t.Fatalf("active connection count = %d, want 0", activeConnections)
	}
	if memberCount != 1 || grp.c[0] != existing {
		t.Fatal("remaining client was removed when another client disconnected")
	}
}

func TestHandleClientReleasesConnectionAfterInitializationFailure(t *testing.T) {
	s := testServer()
	if !s.reserveConnection() {
		t.Fatal("could not reserve a client connection")
	}
	s.wg.Add(1)
	c, peer, _ := testPipeClient(t, s, "")
	peer.Close()

	finished := make(chan struct{})
	go func() {
		handleClient(c)
		close(finished)
	}()
	select {
	case <-finished:
	case <-time.After(2 * time.Second):
		t.Fatal("client handler did not finish after initialization failed")
	}

	s.mu.Lock()
	activeConnections := s.activeConnections
	s.mu.Unlock()
	if activeConnections != 0 {
		t.Fatalf("active connection count = %d after failed initialization, want 0", activeConnections)
	}
}

func TestRunServerReturnsListenError(t *testing.T) {
	if err := runServer(opts{addr: "not a valid listen address"}); err == nil {
		t.Fatal("runServer succeeded with an invalid listen address")
	}
}

func TestAcceptConnectionReturnsErrorAfterListenerCloses(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	if err := listener.Close(); err != nil {
		t.Fatal(err)
	}

	if _, err := acceptConnection(&server{sock: listener}); err == nil {
		t.Fatal("acceptConnection succeeded after the listener closed")
	}
}
