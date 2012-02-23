package httpClient

import(
	"net"
	"url"
	"fmt"
)

const(
	Empty_line = "\r\n"
)

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

func resolveAddr(ip, path *string) *net.TCPAddr{
	host := ip
	if host == nil {
		url,_ := url.Parse(*path)
		host = &url.Host
	}
	addr,_ := net.ResolveTCPAddr("tcp4", *host)
	return addr
}

func do(ip, path, method string, headerMap, parameterMap *map[string]string)string {
	req := method+" "+path+"\r\n"
	req += headers(headerMap)
	req += postBody(parameterMap)
	fmt.Println(req)
	addr := resolveAddr(&ip, &path)
	tcpConn, err := net.DialTCP("tcp4", nil, addr)
	defer tcpConn.Close()
	if err == nil{
		sendRequest(tcpConn, req)
		return recv(tcpConn)
	}
	return ""
}

func Post(ip, path string, headerMap, parameterMap *map[string]string)string{
	return do(ip, path, "POST", headerMap, parameterMap)
}

func Get(ip, path string, headerMap *map[string]string)string{
	return do(ip, path, "GET", headerMap, nil)
}