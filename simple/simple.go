package redis

// import (
// "fmt"
// "net"
// )

// const (
// SIMPLE_STRING_PREFIX = '+'
// ERROR_PREFIX         = '-'
// INTEGER_PREFIX       = ':'
// BULK_STRING_PREFIX   = '$'
// ARRAY_PREFIX         = '*'
// )

// var (
// ConnectionClosed = fmt.Errorf("connection closed")
// )

// type Conn struct {
// // TCP connection to Redis
// conn net.Conn
// // Input buffer
// buf []byte
// // Current index of parser in buffer
// i int
// // Current write position for new data in buffer
// w int
// // Fatal errors
// err error
// }

// func Dial(address string) (*Conn, error) {
// tcpConn, err := net.Dial("tcp", address)
// if err != nil {
// return nil, err
// }

// return &Conn{
// conn: tcpConn,
// }, nil
// }

// func (c *Conn) Do(args ...string) (response []byte, err error) {
// err = c.writeCommand(args...)
// if err != nil {
// return []byte{}, err
// }
// return c.readReply()
// }

// func (c *Conn) writeCommand(args ...string) (err error) {
// // Write the array prefix and the number of arguments to follow.
// _, err = fmt.Fprintf(c.conn, "*%d\r\n", len(args))
// if err != nil {
// return err
// }

// // Send a bulk string for each argument.
// for _, arg := range args {
// _, err = fmt.Fprintf(c.conn, "$%d\r\n%s\r\n", len(arg), arg)
// if err != nil {
// return err
// }
// }

// // All arguments sent successfully!
// return nil
// }

// func (c *Conn) readReply() ([]byte, error) {

// }
