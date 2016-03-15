package main

import (
		"fmt"
		"time"
		"strings"
		"github.com/mxk/go-imap/imap"
)

const(
	Addr = "imap.gmail.com:993"
	User = "test@gmail.com"
	Pass = "123455678"
	MBox = "IMAPBox"
)

func main() {
	c := Dial(Addr)
	defer func(){
		c.Logout(30 * time.Second)
	}()

	c.Noop()
	Login(c, User, Pass)
	cmd, _ := c.List("", "")
	delim :=cmd.Data[0].MailboxInfo().Delim

	mbox := MBox + delim + "Demo1"
	if cmd, err := imap.Wait(c.Create(mbox)); err != nil {
		if rsp, ok := err.(imap.ResponseError); ok && rsp.Status == imap.NO {
			c.Delete(mbox)
		}
		c.Create(mbox)
	} else {
		fmt.Println("error: " , cmd, err)
	}

	fmt.Println("box: " , mbox)
}

func Dial(addr string) (c *imap.Client) {
	var err error
	if strings.HasSuffix(addr, ":993") {
		c, err = imap.DialTLS(addr, nil)
	} else {
		c, err = imap.Dial(addr)
	}
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return c
}

func Login(c *imap.Client, user, pass string) (cmd *imap.Command, err error) {
	defer c.SetLogMask(Sensitive(c, "LOGIN"))
	return c.Login(user, pass)
}

func Sensitive(c *imap.Client, action string) imap.LogMask {
	mask := c.SetLogMask(imap.LogConn)
	hide := imap.LogCmd | imap.LogRaw
	if mask&hide != 0 {
		c.Logln(imap.LogConn, "Raw logging disabled during", action)
	}
	c.SetLogMask(mask &^ hide)
	return mask
}
