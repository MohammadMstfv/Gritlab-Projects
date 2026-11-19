package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
)

const maxClients = 10

type client struct {
	s       *server
	g       *group
	fd      net.Conn
	reader  *bufio.Reader
	name    string
	writeMu sync.Mutex
}

type group struct {
	c       []*client
	name    string
	history string
}

type server struct {
	wg                sync.WaitGroup
	mu                sync.Mutex
	activeConnections int

	g    map[string]*group
	sock net.Listener
	bq   chan message
}

func (s *server) reserveConnection() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.activeConnections >= maxClients {
		return false
	}

	s.activeConnections++
	return true
}

func (s *server) releaseConnection() {
	s.mu.Lock()
	s.activeConnections--
	s.mu.Unlock()
}

func acceptConnection(s *server) (*client, error) {
	conn, e := s.sock.Accept()
	if e != nil {
		return nil, e
	}

	c := client{s: s, fd: conn, reader: bufio.NewReader(conn)}

	return &c, nil
}

func (c *client) write(text string) error {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()

	_, err := io.WriteString(c.fd, text)
	return err
}

func removeClientFromGroupLocked(s *server, c *client) {
	grp := c.g
	if grp == nil {
		return
	}

	for i, member := range grp.c {
		if member == c {
			grp.c = append(grp.c[:i], grp.c[i+1:]...)
			break
		}
	}

	c.g = nil
	if len(grp.c) == 0 {
		delete(s.g, grp.name)
	}
}

func clientGroupAdd(c *client, name string) error {
	s := c.s

	c.writeMu.Lock()
	s.mu.Lock()
	grp, exists := s.g[name]
	if !exists {
		grp = &group{name: name}
		s.g[name] = grp
	}

	if len(grp.c) >= maxClients {
		s.mu.Unlock()
		c.writeMu.Unlock()

		if err := c.write("Chat is full. Try again later.\n"); err != nil {
			return fmt.Errorf("send full-chat message: %w", err)
		}

		return fmt.Errorf("group %q is full", name)
	}

	grp.c = append(grp.c, c)
	c.g = grp
	history := grp.history
	s.mu.Unlock()

	_, err := io.WriteString(c.fd, history)
	if err != nil {
		s.mu.Lock()
		removeClientFromGroupLocked(s, c)
		s.mu.Unlock()
		c.writeMu.Unlock()
		return fmt.Errorf("send group history: %w", err)
	}
	c.writeMu.Unlock()

	broadcastSystemMessage(c, colors["yellow"]+"["+c.name+"]"+colors["reset"]+" joined the group.")

	return nil
}

func clientGroupDel(c *client) error {
	broadcastSystemMessage(c, colors["yellow"]+"["+c.name+"]"+colors["reset"]+" left the group.")
	c.s.mu.Lock()
	defer c.s.mu.Unlock()
	removeClientFromGroupLocked(c.s, c)
	return nil
}

func handleDisconnect(s *server, c *client) error {
	if c.g == nil {
		return nil
	}

	if err := clientGroupDel(c); err != nil {
		return err
	}

	return nil
}

func setName(c *client, name string) error {
	name = strings.TrimSpace(name)

	if name == "" {
		return fmt.Errorf("Names must not be empty\n")
	}
	if len(strings.Fields(name)) > 1 {
		return fmt.Errorf("Names must not include spaces\n")
	}

	if c.name == name {
		return fmt.Errorf("New name identical to the current name\n")
	}

	c.s.mu.Lock()
	defer c.s.mu.Unlock()
	for _, g := range c.s.g {
		for _, cc := range g.c {
			if cc.name == name {
				return fmt.Errorf("This name is already taken\n")
			}
		}
	}

	c.name = name
	return nil
}

func initClient(c *client) error {
	banner := colors["yellow"] + `Welcome to TCP-Chat!
` + colors["green"] + `         _nnnn_
        dGGGGMMb
       @p~qp~~qMb
       M|@||@) M|
       @,----.JM|
      JS^\__/  qKL
     dZP        qKRb
    dZP          qKKb
   fZP            SMMb
   HZM            MMMM
   FqM            MMMM
 __| ".        |\dS"qML
 |    ` + "`" + `.       | ` + "`" + `' \Zq
_)      \.___.,|     .'
\____   )MMMMMP|   .'
     ` + "`" + `-'       ` + "`" + `--'
` + colors["reset"]

	if c.reader == nil {
		c.reader = bufio.NewReader(c.fd)
	}

	if err := c.write(banner); err != nil {
		return fmt.Errorf("send welcome banner: %w", err)
	}

	for {
		if err := c.write(colors["yellow"] + "[ENTER YOUR NAME]: " + colors["reset"]); err != nil {
			return fmt.Errorf("send name prompt: %w", err)
		}

		name, err := c.reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("read client name: %w", err)
		}

		if err = setName(c, name); err != nil {
			if err := c.write(err.Error()); err != nil {
				return fmt.Errorf("send name guidance: %w", err)
			}
			continue
		}

		if err := clientGroupAdd(c, ""); err != nil {
			return fmt.Errorf("add client to default group: %w", err)
		}

		return nil
	}
}

func clientLoop(c *client) error {
	for {
		msg, err := receiveMessage(c)
		if err != nil {
			return err
		}
		if err := processMessage(msg); err != nil {
			return err
		}
	}
}

func handleClient(c *client) {
	defer c.s.wg.Done()
	defer c.s.releaseConnection()
	defer c.fd.Close()

	defer handleDisconnect(c.s, c)

	if err := initClient(c); err != nil {
		log.Printf("client %s initialization failed: %v", c.fd.RemoteAddr(), err)
		return
	}

	if err := clientLoop(c); err != nil {
		log.Printf("client %s disconnected: %v", c.fd.RemoteAddr(), err)
	}
}

func broadcastLoop(s *server) {
	defer s.wg.Done()

	for {
		msg := <-s.bq
		if err := broadcastMessage(msg); err != nil {
			log.Printf("broadcast chat message: %v", err)
		}
	}
}

func connectLoop(s *server) {
	defer s.wg.Done()

	for {
		c, err := acceptConnection(s)
		if err != nil {
			log.Printf("accept connection: %v", err)
			continue
		}

		if !s.reserveConnection() {
			if _, err := io.WriteString(c.fd, "Chat is full. Try again later.\n"); err != nil {
				log.Printf("reject client %s: %v", c.fd.RemoteAddr(), err)
			}
			if err := c.fd.Close(); err != nil {
				log.Printf("close rejected client %s: %v", c.fd.RemoteAddr(), err)
			}
			continue
		}

		s.wg.Add(1)
		go handleClient(c)
	}
}

func runServer(op opts) error {
	sock, e := net.Listen("tcp", op.addr)
	if e != nil {
		return e
	}
	defer sock.Close()

	fmt.Printf("Listening on the port %s\n", op.addr)

	s := server{
		g:    make(map[string]*group),
		sock: sock,
		bq:   make(chan message),
	}

	s.wg.Add(2)
	go broadcastLoop(&s)
	go connectLoop(&s)

	s.wg.Wait()

	return nil
}
