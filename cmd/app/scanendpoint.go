// +build !windows

package app

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/anchnet/smartops-agent/pkg/http"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/spf13/cobra"
)

var (
	scanendpointCmd = &cobra.Command{
		Use:   "scanendpoint",
		Short: "Stand scanendpoint",
		RunE:  scanendpoint,
	}
	ScanData   map[string]Info
	localIface string
	LocalIpNet *net.IPNet
	localMac   net.HardwareAddr

	t  *time.Ticker //计时器，在一段时间没新数据流入， 推出程序
	do chan string
)

type IP uint32

// 将 IP(uint32) 转换成 可读性IP字符串
func (ip IP) String() string {
	var bf bytes.Buffer
	for i := 1; i <= 4; i++ {
		bf.WriteString(strconv.Itoa(int((ip >> ((4 - uint(i)) * 8)) & 0xff)))
		if i != 4 {
			bf.WriteByte('.')
		}
	}
	return bf.String()
}

type Info struct {
	// MAC地址
	Mac string `json:"mac"`
	// 主机名
	Hostname string `json:"hostname"`
	// IP地址
	IP string `json:"privateIP"`
	//time
	Time string `json:"time"`
}

func ScanDataTOJson() ([]byte, error) {
	var data []Info
	data = make([]Info, 0)
	for ip, info := range ScanData {
		hostname := info.Hostname
		if hostname == "" {
			hostname = "Unknown"
		}
		info := Info{
			Mac:      info.Mac,
			Hostname: hostname,
			IP:       ip,
			Time:     time.Now().Format("2006-01-02T15:04:05.999+07:00"),
		}
		data = append(data, info)
	}
	return json.Marshal(data)
}

func scanendpoint(cmd *cobra.Command, args []string) error {

	do = make(chan string)
	ScanData = make(map[string]Info)

	// 获取本地信息
	ipNet, iface, mac, err := getLocalIP()
	if err != nil {
		return err
	}
	localIface = iface
	LocalIpNet = ipNet
	localMac = mac
	host, _ := os.Hostname()
	fmt.Printf("本地IP:%s 本地hostname:%s, 本地mac: %s \n", ipNet.IP.String(), host, localMac.String())

	ctx, cancel := context.WithCancel(context.Background())
	go listenARP(ctx)
	go sendARP(ipNet)
	go listenMDNS(ctx)

	t = time.NewTicker(20 * time.Second)
	for {
		select {
		case <-t.C:
			cancel()
			byts, _ := ScanDataTOJson()
			fmt.Println(string(byts))
			respByte, err := http.PhysicalDevice(byts)
			if err != nil {
				fmt.Println("Send local metric server error: ", err)
				return err
			}
			fmt.Println(respByte)
			goto END
		case d := <-do:
			switch d {
			case "start":
				t.Stop()
			case "end":
				// 接收到新数据，重置2秒的计数器
				t = time.NewTicker(20 * time.Second)
			}
		}
	}
END:
	fmt.Println("完成扫描")
	return err
}

func sendARP(ipNet *net.IPNet) {
	ips, min, max := Table(ipNet)
	fmt.Printf("扫描的范围为: %s ---- %s\n", min, max)
	for _, ip := range ips {
		time.Sleep(time.Millisecond * 20)
		go sendArpPackage(ip)
	}
}

func pushData(ip string, mac net.HardwareAddr, hostname string) {
	// 停止计时器
	do <- "start"
	var mu sync.RWMutex
	mu.RLock()
	defer func() {
		// 重置计时器
		do <- "end"
		mu.RUnlock()
	}()
	if _, ok := ScanData[ip]; !ok {
		ScanData[ip] = Info{Mac: mac.String(), Hostname: hostname}
		return
	}
	info := ScanData[ip]
	if len(hostname) > 0 && len(info.Hostname) == 0 {
		info.Hostname = hostname
	}

	if mac != nil {
		info.Mac = mac.String()
	}

	ScanData[ip] = info
}

func listenARP(ctx context.Context) {
	handle, err := pcap.OpenLive(localIface, 1024, false, 10*time.Second)
	if err != nil {
		log.Fatal("监听ARP pcap打开失败:", err)
	}
	fmt.Println("开始监听ARP")
	defer handle.Close()
	handle.SetBPFFilter("arp")
	ps := gopacket.NewPacketSource(handle, handle.LinkType())
	for {
		select {
		case <-ctx.Done():
			return
		case p := <-ps.Packets():
			arp := p.Layer(layers.LayerTypeARP).(*layers.ARP)
			if arp.Operation == 2 {
				mac := net.HardwareAddr(arp.SourceHwAddress)
				pushData(ParseIP(arp.SourceProtAddress).String(), mac, "")
				go sendMdns(ParseIP(arp.SourceProtAddress), mac)
			}
		}
	}
}

// []byte --> IP
func ParseIP(b []byte) IP {
	return IP(IP(b[0])<<24 + IP(b[1])<<16 + IP(b[2])<<8 + IP(b[3]))
}

func sendArpPackage(ip IP) error {
	srcIp := net.ParseIP(LocalIpNet.IP.String()).To4()
	dstIp := net.ParseIP(ip.String()).To4()
	if srcIp == nil || dstIp == nil {
		return errors.New("ip 解析出问题")
	}
	// 以太网首部
	// EthernetType 0x0806  ARP
	ether := &layers.Ethernet{
		SrcMAC:       localMac,
		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		EthernetType: layers.EthernetTypeARP,
	}

	a := &layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		HwAddressSize:     uint8(6),
		ProtAddressSize:   uint8(4),
		Operation:         uint16(1), // 0x0001 arp request 0x0002 arp response
		SourceHwAddress:   localMac,
		SourceProtAddress: srcIp,
		DstHwAddress:      net.HardwareAddr{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		DstProtAddress:    dstIp,
	}

	buffer := gopacket.NewSerializeBuffer()
	var opt gopacket.SerializeOptions
	gopacket.SerializeLayers(buffer, opt, ether, a)
	outgoingPacket := buffer.Bytes()

	handle, err := pcap.OpenLive(localIface, 2048, false, 30*time.Second)
	if err != nil {
		return errors.New("发送ARP pcap打开失败:" + err.Error())
	}
	defer handle.Close()

	err = handle.WritePacketData(outgoingPacket)
	if err != nil {
		return errors.New("发送arp数据包失败")
	}
	return nil
}

// 根据IP和mask换算内网IP范围
func Table(ipNet *net.IPNet) ([]IP, IP, IP) {
	ip := ipNet.IP.To4()
	var min, max IP
	var data []IP
	for i := 0; i < 4; i++ {
		b := IP(ip[i] & ipNet.Mask[i])
		min += b << ((3 - uint(i)) * 8)
	}
	one, _ := ipNet.Mask.Size()
	max = min | IP(math.Pow(2, float64(32-one))-1)
	// max 是广播地址，忽略
	// i & 0x000000ff  == 0 是尾段为0的IP，根据RFC的规定，忽略
	for i := min; i < max; i++ {
		if i&0x000000ff == 0 {
			continue
		}
		data = append(data, i)
	}
	return data, min, max
}

func getLocalIP() (ipNet *net.IPNet, iface string, mac net.HardwareAddr, err error) {
	ifs, err := net.Interfaces()
	if err != nil {
		return
	}
	for _, it := range ifs {
		addr, _ := it.Addrs()
		for _, a := range addr {
			if ip, ok := a.(*net.IPNet); ok && !ip.IP.IsLoopback() {
				if ip.IP.To4() != nil {
					ipNet = ip
					mac = it.HardwareAddr
					iface = it.Name
					goto END
				}
			}
		}
	}

END:
	if ipNet == nil || len(mac) == 0 {
		err = errors.New("无法获取本地网络信息")
	}
	return
}

func init() {
	Command.AddCommand(scanendpointCmd)
}

//MDNS ==================================================================
func listenMDNS(ctx context.Context) {
	handle, err := pcap.OpenLive(localIface, 1024, false, 10*time.Second)
	if err != nil {
		log.Fatal("linsten MDNS pcap打开失败:", err)
	}
	defer handle.Close()
	handle.SetBPFFilter("udp and port 5353")
	ps := gopacket.NewPacketSource(handle, handle.LinkType())
	for {
		select {
		case <-ctx.Done():
			return
		case p := <-ps.Packets():
			if len(p.Layers()) == 4 {
				c := p.Layers()[3].LayerContents()
				if c[2] == 0x84 && c[3] == 0x00 && c[6] == 0x00 && c[7] == 0x01 {
					// 从网络层(ipv4)拿IP, 不考虑IPv6
					i := p.Layer(layers.LayerTypeIPv4)
					if i == nil {
						continue
					}
					ipv4 := i.(*layers.IPv4)
					ip := ipv4.SrcIP.String()
					// 把 hostname 存入到数据库
					h := ParseMdns(c)
					if len(h) > 0 {
						pushData(ip, nil, h)
					}
				}
			}
		}
	}
}

type Buffer struct {
	data  []byte
	start int
}

func (b *Buffer) PrependBytes(n int) []byte {
	length := cap(b.data) + n
	newData := make([]byte, length)
	copy(newData, b.data)
	b.start = cap(b.data)
	b.data = newData
	return b.data[b.start:]
}

func NewBuffer() *Buffer {
	return &Buffer{}
}

// 反转字符串
func Reverse(s string) (result string) {
	for _, v := range s {
		result = string(v) + result
	}
	return
}

// 根据ip生成含mdns请求包，包存储在 buffer里
func mdns(buffer *Buffer, ip string) {
	b := buffer.PrependBytes(12)
	binary.BigEndian.PutUint16(b, uint16(0))          // 0x0000 标识
	binary.BigEndian.PutUint16(b[2:], uint16(0x0100)) // 标识
	binary.BigEndian.PutUint16(b[4:], uint16(1))      // 问题数
	binary.BigEndian.PutUint16(b[6:], uint16(0))      // 资源数
	binary.BigEndian.PutUint16(b[8:], uint16(0))      // 授权资源记录数
	binary.BigEndian.PutUint16(b[10:], uint16(0))     // 额外资源记录数
	// 查询问题
	ipList := strings.Split(ip, ".")
	for j := len(ipList) - 1; j >= 0; j-- {
		ip := ipList[j]
		b = buffer.PrependBytes(len(ip) + 1)
		b[0] = uint8(len(ip))
		for i := 0; i < len(ip); i++ {
			b[i+1] = uint8(ip[i])
		}
	}
	b = buffer.PrependBytes(8)
	b[0] = 7 // 后续总字节
	copy(b[1:], []byte{'i', 'n', '-', 'a', 'd', 'd', 'r'})
	b = buffer.PrependBytes(5)
	b[0] = 4 // 后续总字节
	copy(b[1:], []byte{'a', 'r', 'p', 'a'})
	b = buffer.PrependBytes(1)
	// terminator
	b[0] = 0
	// type 和 classIn
	b = buffer.PrependBytes(4)
	binary.BigEndian.PutUint16(b, uint16(12))
	binary.BigEndian.PutUint16(b[2:], 1)
}

func sendMdns(ip IP, mhaddr net.HardwareAddr) {
	srcIp := net.ParseIP(LocalIpNet.IP.String()).To4()
	dstIp := net.ParseIP(ip.String()).To4()
	ether := &layers.Ethernet{
		SrcMAC:       localMac,
		DstMAC:       mhaddr,
		EthernetType: layers.EthernetTypeIPv4,
	}

	ip4 := &layers.IPv4{
		Version:  uint8(4),
		IHL:      uint8(5),
		TTL:      uint8(255),
		Protocol: layers.IPProtocolUDP,
		SrcIP:    srcIp,
		DstIP:    dstIp,
	}
	bf := NewBuffer()
	mdns(bf, ip.String())
	udpPayload := bf.data
	udp := &layers.UDP{
		SrcPort: layers.UDPPort(60666),
		DstPort: layers.UDPPort(5353),
	}
	udp.SetNetworkLayerForChecksum(ip4)
	udp.Payload = udpPayload // todo
	buffer := gopacket.NewSerializeBuffer()
	opt := gopacket.SerializeOptions{
		FixLengths:       true, // 自动计算长度
		ComputeChecksums: true, // 自动计算checksum
	}
	err := gopacket.SerializeLayers(buffer, opt, ether, ip4, udp, gopacket.Payload(udpPayload))
	if err != nil {
		log.Fatal("Serialize layers出现问题:", err)
	}
	outgoingPacket := buffer.Bytes()

	handle, err := pcap.OpenLive(localIface, 1024, false, 10*time.Second)
	if err != nil {
		log.Fatal("发送MDNS pcap打开失败:", err)
	}
	defer handle.Close()
	err = handle.WritePacketData(outgoingPacket)
	if err != nil {
		log.Fatal("发送udp数据包失败..")
	}
}

// 参数data  开头是 dns的协议头 0x0000 0x8400 0x0000 0x0001(ans) 0x0000 0x0000
// 从 mdns响应报文中获取主机名
func ParseMdns(data []byte) string {
	var buf bytes.Buffer
	i := bytes.Index(data, []byte{0x05, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x00})
	if i < 0 {
		return ""
	}

	for s := i - 1; s > 1; s-- {
		num := i - s
		if s-2 < 0 {
			break
		}
		// 包括 .local_ 7 个字符
		if bto16([]byte{data[s-2], data[s-1]}) == uint16(num+7) {
			return Reverse(buf.String())
		}
		buf.WriteByte(data[s])
	}

	return ""
}

func bto16(b []byte) uint16 {
	if len(b) != 2 {
		log.Fatal("b只能是2个字节")
	}
	return uint16(b[0])<<8 + uint16(b[1])
}
