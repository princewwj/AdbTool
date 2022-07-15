package main
 
import (
    "fmt"
    "net"
    "os"
    "strconv"
    "strings"
    "time"
)

const WILL = 251
const WONT = 252
const DO = 253
const DONT = 254
const IAC = 255
const RD = 1
const SGA = 3
 
func main() {

    var n int=0
    var count int=0

    var buf [8196]byte
    var kkk[100]byte

    var srcIP string="0.0.0.0"
    var timeout int=16

    var index int=1
    var params[68] string

    for k:=1;k<len(os.Args);k++{
        if os.Args[k]=="--interface" || os.Args[k]=="-S" {
            k++
            srcIP=os.Args[k]
        }else if os.Args[k]=="--timeout" || os.Args[k]=="-t" {
            k++
            timeout,_=strconv.Atoi(os.Args[k])
        }else{
            params[index]=os.Args[k]
            index++
        }
    }

    var localaddr net.TCPAddr
    var remoteaddr net.TCPAddr

    localaddr.IP = net.ParseIP(srcIP)
    localaddr.Port = 0

    remoteaddr.IP = net.ParseIP(params[1])
    remoteaddr.Port,_= strconv.Atoi(params[2])

    conn, err := net.DialTCP("tcp", &localaddr, &remoteaddr)
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(-1)
    }

    conn.SetReadDeadline(time.Now().Add(time.Second * 8))
    for {
        n, err = conn.Read(buf[0:])
        count:=n/3
        if buf[n-1-1]==byte(':') {
            fmt.Println(string(buf[0:n]))
            break
        }
        for i:=0;i<count;i++ {
            iac := buf[i * 3]   
            cmd := buf[i * 3 + 1]
            value := buf[i * 3 + 2]
            if (IAC != iac) {
                continue;
            }
            switch (cmd){
                case DO:
                    kkk[0]=byte(iac)
                    if value==RD {
                        kkk[1]=byte(WILL)
                    }else{
                        kkk[1]=byte(WONT)
                    }
                    kkk[2]=byte(value)
                    n, err = conn.Write(kkk[0:3])
                    break
                case DONT:
                    kkk[0]=byte(iac)
                    kkk[1]=byte(WONT)
                    kkk[2]=byte(value)
                    n, err = conn.Write(kkk[0:3])
                    break
                case WILL:
                    kkk[0]=byte(iac)
                    if value==SGA {
                    kkk[1]=byte(DO)
                    }else{
                        kkk[1]=byte(DONT)
                    }
                    kkk[2]=byte(value)
                    n, err = conn.Write(kkk[0:3])
                    break
                case WONT:
                    kkk[0]=byte(iac)
                    kkk[1]=byte(DONT)
                    kkk[2]=byte(value)
                    n, err = conn.Write(kkk[0:3])
                    break
                 default:
                     break
             }
        }
    }

    if strings.Contains(string(buf[0:n]),"ogin:") {
        n, err = conn.Write([]byte(params[3]+"\r\n"))
        if err != nil {
            fmt.Println(err.Error())
            os.Exit(-1)
        }    
    }else{
        fmt.Println("Error!")
        os.Exit(-1)
    }

    count=0    
    conn.SetReadDeadline(time.Now().Add(time.Second * 8))
    for k:=0;k<len(buf);k++{
        buf[k]=0
    }

    for k:=0;k<68;k++ {
        n, err = conn.Read(buf[count:])
        if err != nil {
            fmt.Println(err.Error())
            break
        }
        count += n
        if strings.Contains(string(buf[0:count]),"assword:") {
            break
        }
        time.Sleep(time.Microsecond * 100000)
    }

    fmt.Print(string(buf[0:count]))
    if strings.Contains(string(buf[0:count]),"assword:") {
        n, err = conn.Write([]byte(params[3+1]+"\r\n"))
        if err != nil {
            fmt.Println(err.Error())
            os.Exit(-1)
        }
    }else{
        fmt.Println("Error!!!")
        os.Exit(-1)
    }

    conn.SetReadDeadline(time.Now().Add(time.Second * 8))
    for k:=0;k<68;k++ {
        n, err = conn.Read(buf[0:])
        if err != nil {
            fmt.Println(err.Error())
            break
        }
        fmt.Print(string(buf[0:n]))
        if strings.Contains(string(buf[0:n]),"# ") {
            break
        }
        time.Sleep(time.Microsecond * 100000)
    }

    n, err = conn.Write([]byte(params[5]+"\r\n"))
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(-1)
    }

    conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(timeout)))
    for k:=0;k<timeout*100;k++ {
        n, err = conn.Read(buf[0:])
        if err != nil {
            fmt.Println(err.Error())
            break
        }
        fmt.Print(string(buf[0:n]))
        if strings.Contains(string(buf[0:n]),"# ") {
            break
        }
        time.Sleep(time.Microsecond * 10000)
    }
 
    os.Exit(0)
}