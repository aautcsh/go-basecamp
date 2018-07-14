package main

import (
	"bufio"
	"io"
	"log"
	"math"
	"math/big"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	msgdir  = "./msg/"
	koandir = "./koans/"
)

var (
	procedures = map[string]func(*Processor, []string){
		"rand": func(p *Processor, c []string) {
			if len(c) == 0 {
				p.WriteFlush("02 ERROR NO NUMBER SPECIFIED")
				return
			}

			n, err := strconv.ParseUint(c[0], 10, 64)
			if err != nil {
				p.WriteFlush("02 ERROR " + c[0] + " INVALID")
				return
			}
			p.Write("01 OK")
			var i uint64
			for i = 0; i < n; i++ {
				p.Write(strconv.Itoa(int(rand.Int31n(256))))
			}
			p.WriteFlush(".")
		},
		"quine": func(p *Processor, c []string) {
			p.Write("01 OK")
			f, err := os.OpenFile("quine.go", os.O_RDONLY, 0600)
			if err != nil {
				p.WriteFlush("02 ERROR QUINE NOT AVAILABLE")
				return
			}
			defer f.Close()
			io.Copy(p.Bout, f)
			p.WriteFlush(".")
		},
		"base29": func(p *Processor, c []string) {
			if len(c) == 0 {
				p.WriteFlush("02 ERROR NO NUMBER SPECIFIED")
				return
			}

			n, err := strconv.ParseUint(c[0], 10, 64)
			if err != nil {
				p.WriteFlush("02 ERROR " + c[0] + " INVALID")
				return
			}

			s := ""
			for {
				r := n % 29
				s = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"[r:r+1] + s
				n = n / 29
				if n == 0 {
					break
				}
			}
			p.WriteFlush("01 OK " + s)
		},
		"code": func(p *Processor, c []string) {
			p.Write("01 OK")
			f, err := os.OpenFile("c3301serv.go", os.O_RDONLY, 0600)
			if err != nil {
				p.WriteFlush("02 ERROR SOURCE NOT AVAILABLE")
				return
			}
			defer f.Close()
			io.Copy(p.Bout, f)
			p.WriteFlush(".")
		},
		"koan": func(p *Processor, c []string) {
			d, err := os.Open(koandir)
			if err != nil {
				p.WriteFlush("02 ERROR KOAN NOT AVAILABLE")
				return
			}
			defer d.Close()
			fs, err := d.Readdirnames(-1)
			if err != nil {
				p.WriteFlush("02 ERROR KOAN NOT AVAILABLE")
				return
			}
			f, err := os.OpenFile(koandir+fs[rand.Intn(len(fs))], os.O_RDONLY, 0600)
			if err != nil {
				p.WriteFlush("02 ERROR KOAN NOT AVAILABLE")
				return
			}
			defer f.Close()
			p.Write("01 OK")
			io.Copy(p.Bout, f)
			p.WriteFlush(".")
		},
		"dh": func(p *Processor, c []string) {
			if len(c) == 0 {
				p.WriteFlush("02 ERROR NO PRIME MODULUS")
				return
			}

			m, ok := big.NewInt(0).SetString(c[0], 10)
			if !ok {
				p.WriteFlush("02 ERROR " + c[0] + " INVALID")
				return
			}
			b := big.NewInt(rand.Int63n(23) + 3)
			s := big.NewInt(rand.Int63n(math.MaxInt64-3301) + 3301)
			e := big.NewInt(0).Exp(b, s, m)
			p.WriteFlush("01 OK " + b.String() + " " + e.String())

			p.C.SetReadDeadline(time.Now().Add(time.Second * 15))
			l, err := p.ReadLine()
			if err != nil {
				p.WriteFlush("02 ERROR TIMEOUT")
				return
			}
			e2, ok := big.NewInt(0).SetString(l, 10)
			if !ok {
				p.WriteFlush("02 ERROR " + l + " INVALID")
				return
			}
			k := big.NewInt(0).Exp(e2, s, m)
			p.WriteFlush("03 DATA " + k.String())
		},
		"next": func(p *Processor, c []string) {
			fn := strconv.FormatUint(uint64(time.Now().Unix()), 10)
			f, err := os.OpenFile(msgdir+fn, os.O_WRONLY|os.O_CREATE, 0700)
			if err != nil {
				p.WriteFlush("02 ERROR INTERNAL PROBLEM")
				return
			}
			defer f.Close()

			p.WriteFlush("01 OK")
			for {
				p.C.SetReadDeadline(time.Now().Add(time.Second * 15))
				l, err := p.ReadLine()
				if err != nil {
					p.WriteFlush("02 ERROR TIMEOUT")
					return
				}

				if l == "." {
					break
				}

				f.Write([]byte(l))
				f.Write([]byte{'\n'})
			}
			p.WriteFlush("01 OK")
		},
		"goodbye": func(p *Processor, c []string) {
			p.WriteFlush("99 GOODBYE")
			p.Stop()
		},
	}
)

type Processor struct {
	C    *net.TCPConn
	Bout *bufio.Writer
	Bin  *bufio.Reader
	Log  *log.Logger
	R    bool
}

func NewProcessor(c *net.TCPConn, l *log.Logger) *Processor {
	return &Processor{c, bufio.NewWriter(c), bufio.NewReader(c), l, false}
}

func (this *Processor) Start() {
	this.R = true
	this.WriteFlush("00 WELCOME TO CICADA SERVICE WRITTEN IN GO")
	for this.R {
		this.C.SetReadDeadline(time.Now().Add(time.Minute))
		l, err := this.ReadLine()
		if err != nil {
			this.WriteFlush("02 ERROR TIMEOUT")
			this.Stop()
			break
		}

		c := strings.Split(strings.ToLower(l), " ")
		if len(c) == 0 {
			this.WriteFlush("02 ERROR NO COMMAND")
			continue
		}

		p, ok := procedures[c[0]]
		if !ok {
			this.WriteFlush("02 ERROR COMMAND " + c[0] + " INVALID")
			continue
		}
		this.Log.Print("executing '" + strings.Join(c, " ") + "' for " + this.C.RemoteAddr().String())
		this.C.SetReadDeadline(time.Unix(0, 0))
		p(this, c[1:])
	}
}

func (this *Processor) WriteFlush(data string) {
	this.Bout.WriteString(data + "\r\n")
	this.Bout.Flush()
}

func (this *Processor) Write(data string) {
	this.Bout.WriteString(data + "\r\n")
}

func (this *Processor) Flush() {
	this.Bout.Flush()
}

func (this *Processor) Stop() {
	this.Log.Print("closing connection to " + this.C.RemoteAddr().String())
	this.C.Close()
	this.R = false
}

func (this *Processor) ReadLine() (string, error) {
	l, err := this.Bin.ReadString('\n')
	if err != nil {
		return "", err
	}
	if len(l) >= 2 {
		l = l[:len(l)-2]
	}
	return l, nil
}

type Server struct {
	L   *net.TCPListener
	Log *log.Logger
}

func NewServer(laddr *net.TCPAddr) (*Server, error) {
	l, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		return nil, err
	}
	return &Server{l, log.New(os.Stdout, "", log.LstdFlags)}, nil
}

func (this *Server) Run() {
	this.Log.Println("server is running under address " + this.L.Addr().String())
	for {
		c, err := this.L.AcceptTCP()
		if err != nil {
			this.Log.Print(err)
			continue
		}

		raddr := c.RemoteAddr().(*net.TCPAddr)
		if !raddr.IP.IsLoopback() {
			c.Close()
			continue
		}

		this.Log.Println("got connection from " + c.RemoteAddr().String())
		p := NewProcessor(c, this.Log)
		go p.Start()
	}
}

func main() {
	laddr, err := net.ResolveTCPAddr("tcp", ":3307")
	if err != nil {
		log.Fatal(err)
	}
	s, err := NewServer(laddr)
	if err != nil {
		log.Fatal(err)
	}
	s.Run()
}
