package main

import "testing"

func TestTransform(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "case markers with counts",
			in:   "it (cap) was the best of times, it was the worst of times (up) , it was the age of wisdom, it was the age of foolishness (cap, 6) , it was the epoch of belief, it was the epoch of incredulity, it was the season of Light, it was the season of darkness, it was the spring of hope, IT WAS THE (low, 3) winter of despair.",
			want: "It was the best of times, it was the worst of TIMES, it was the age of wisdom, It Was The Age Of Foolishness, it was the epoch of belief, it was the epoch of incredulity, it was the season of Light, it was the season of darkness, it was the spring of hope, it was the winter of despair.",
		},
		{
			name: "hex and bin",
			in:   "Simply add 42 (hex) and 10 (bin) and you will see the result is 68.",
			want: "Simply add 66 and 2 and you will see the result is 68.",
		},
		{
			name: "a to an before vowel",
			in:   "There is no greater agony than bearing a untold story inside you.",
			want: "There is no greater agony than bearing an untold story inside you.",
		},
		{
			name: "punctuation spacing and dot groups",
			in:   "Punctuation tests are ... kinda boring ,what do you think ?",
			want: "Punctuation tests are... kinda boring, what do you think?",
		},
		{
			name: "single word quotes",
			in:   "I am exactly how they describe me: ' awesome '",
			want: "I am exactly how they describe me: 'awesome'",
		},
		{
			name: "multi word quotes",
			in:   "As Elton John said: ' I am the most well-known homosexual in the world '",
			want: "As Elton John said: 'I am the most well-known homosexual in the world'",
		},
		{
			name: "capital A to An",
			in:   "There it was. A amazing rock!",
			want: "There it was. An amazing rock!",
		},
		// Extra edge cases not directly in the spec examples.
		{
			name: "a to an before h",
			in:   "He was a hero.",
			want: "He was an hero.",
		},
		{
			name: "up with no count defaults to one word",
			in:   "make this LOUD (up)",
			want: "make this LOUD",
		},
		{
			name: "marker at start of line has nothing to act on",
			in:   "(cap) hello there",
			want: "hello there",
		},
		{
			name: "exclamation and question mark group stays glued",
			in:   "Wait, what !? Really.",
			want: "Wait, what!? Really.",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := Transform(c.in)
			if got != c.want {
				t.Errorf("Transform(%q)\n got:  %q\n want: %q", c.in, got, c.want)
			}
		})
	}
}
