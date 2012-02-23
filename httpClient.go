package httpClient

import(
	"net"
	"url"
	"fmt"
	"strings"
)

const(
	Empty_line = "\r\n"
)

type HttpClient struct{
	connectTimeout int
	readTimeout int
	writeTimeout int
	remoteIp string
}

func postBody(parameterMap *map[string]string)string{
	postBody := ""
	if parameterMap != nil {
		for pn, pv := range *parameterMap {
			postBody += pn+"="+pv+"&"
		}
		postBody += fmt.Sprintf("Content-Length: %d\r\n", len(postBody))
		postBody += Empty_line
	}
	return postBody
}

func headers(headerMap *map[string]string)string{
	result := ""
	if headerMap != nil {
		for hn, hv := range *headerMap {
			result += hn + ": " + hv + "\r\n"
		}
	}
	return result
}

func sendRequest(tcpConn *net.TCPConn, req string){
	reqBytes := []byte(req)
	total := len(reqBytes)
	n, err := tcpConn.Write(reqBytes)
	offset := n
	for err == nil && offset < total{
		n, err = tcpConn.Write(reqBytes[offset:total])
		offset += n
	}
}

func recv(tcpConn *net.TCPConn)string{
	result := ""
	buffer := make([]byte, 1024)
	n, err := tcpConn.Read(buffer)
	result = string(buffer[0:n])
	for err == nil {
		n, err = tcpConn.Read(buffer)	
		result += string(buffer[0:n])
	}
	return result
}

func appendPort(host string)string{
	result := ""
	index := strings.Index(host, ":")
	if index == -1{
		result = host + ":80"
	}
	if index == (len(host)-1) {
		result = host + "80"
	}
	return result
}

func resolveAddr(ip, path string) *net.TCPAddr{
	host := ip
	if  len(host) == 0{
		url,_ := url.Parse(path)
		host = url.Host
	}
	host = appendPort(host)
	addr,_ := net.ResolveTCPAddr("tcp4", host)
	return addr
}

func (c *HttpClient)SetRemoteIp(ip string){
	c.remoteIp = ip
}

func (c *HttpClient)do(url, method string, headerMap, parameterMap *map[string]string)string {
	req := method+" "+url+"\r\n"
	req += headers(headerMap)
	req += postBody(parameterMap)
	addr := resolveAddr(c.remoteIp, url)
	tcpConn, err := net.DialTCP("tcp4", nil, addr)
	defer tcpConn.Close()
	if err == nil{
		sendRequest(tcpConn, req)
		return recv(tcpConn)
	}
	return ""
}

func (c *HttpClient)Post(url string, headerMap, parameterMap *map[string]string)string{
	return c.do(url, "POST", headerMap, parameterMap)
}

func (c *HttpClient)Get(url string, headerMap *map[string]string)string{
	return c.do(url, "GET", headerMap, nil)
}

func CreateHttpClient(connectTimeout, readTimeout, writeTimeout int)(*HttpClient){
	return &HttpClient{connectTimeout:connectTimeout, readTimeout:readTimeout, writeTimeout:writeTimeout}
}