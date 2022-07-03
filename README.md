# mbox

Supporting four different mbox file formats in pure golang.

Package mbox implements a reader and writer for working with mbox files.

The package supports four types of mbox files:

- mboxo
- mboxrd
- mboxcl
- mboxcl2

Type `mboxo` is the original mbox format.

Type `mboxrd` tries to address lines starting with 'From ' in a way to avoid
conflicts by prepending such lines with '>', removing those characters when
reading the mail.

Type `mboxcl` tries to address lines starting with 'From ' by doing what mboxrd
does, but also adding a 'Content-Length' header to the mail that provides the
size of the mail's body.

Type `mboxcl2` tries to address the lines starting with 'From ' by doing what
mboxcl does, except it doesn't add '>' characters at all.

You will need to know which type to use when reading or writing an mbox, for
best results.

NOTE: These routines do not concern themselves with file locking. You may want
to consider that while working with mbox files on systems that might actively
write to the file. These simply use the golang writer/reader interfaces.

General usage:

Reading mbox files:

```golang
    import (
     "bytes"
     "fmt"
     "io"
     "net/mail"
     "os"

     "github.com/tvanriper/mbox"
    )

 file, err := os.Open(filepath)
 if err != nil {
  panic(err)
 }
 box := mbox.NewReader(file)
 box.Type = mbox.MBOXRD // chooses mbox.MBOXO by default

 var from string
 // NextMessage will return err == io.EOF when the last message has been read.
 for err != io.EOF {
  data := bytes.NewBuffer([]byte{})
  from, err = box.NextMessage(data)
  if err != nil && err != io.EOF {
   panic(err)
  }
  // You normally do not need this, but it's provided in case you do.
  fmt.Printf("from line: %s\n", from)

  // You can use net/mail to parse the mail itself
  msg, e := mail.ReadMessage(bytes.NewBuffer(data.Bytes()))
  if e != nil {
   fmt.Printf("problems reading mail: %s\n", e)
   continue
  }

  fmt.Printf("From: %s\n", msg.Header.Get("From"))
  fmt.Printf("Subject: %s\n", msg.Header.Get("Subject"))
  fmt.Printf("To: %s\n", msg.Header.Get("To"))
  fmt.Println("Body:")
  io.Copy(os.Stdout, msg.Body)
  fmt.Println("")
 }
```

Writing mbox files:

```golang
file, err := os.Create(filepath)
if err != nil {
  panic(err)
}

box := NewWriter(file)
box.Type = mbox.MBOXRD
// 'email' holds both headers and body of the email message you want to add to
// the mbox file in a single []byte.
email := YourMailToBytesFn()
err = mbox.WriteMail(mbox.ParseFrom("bubbahotep@hotelcalifornia.com",time.Now(),""), email)
if err != nil {
  panic(err)
}
email = YourEmailToBytesFn()
err = mbox.WriteMail(mbox.ParseFrom("notintheface@politics.com", time.Now(), ""), email)
if err != nil {
    panic(err)
}
```
