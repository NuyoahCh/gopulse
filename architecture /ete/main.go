package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"
)

// 初始化参数
var (
	sizeMB       = flag.Int("sizeMB", 2, "模拟文件大小（MB）")
	packetSize   = flag.Int("packet", 1024, "分包大小（字节）")
	linkFlipProb = flag.Float64("flipProb", 0.01, "链路每包发生比特翻转的概率 [0,1]")
	diskFlipProb = flag.Float64("diskProb", 0.02, "接收端写盘后发生单字节腐化的概率 [0,1]")
	linkReliable = flag.Bool("linkReliable", true, "是否启用链路层校验+重传")
	maxRetries   = flag.Int("retries", 5, "端到端重试次数上限")
	destPath     = flag.String("out", "dest.bin", "接收端输出文件")
	randomSeed   = flag.Int64("seed", time.Now().UnixNano(), "随机种子")
)

func main() {
	flag.Parse()
	rand.Seed(*randomSeed)

	fmt.Printf("=== End-to-End Demo (MIT E2E Argument) ===\n")
	fmt.Printf("size=%dMB packet=%d flipProb=%.3f diskProb=%.3f linkReliable=%v retries=%d\n\n",
		*sizeMB, *packetSize, *linkFlipProb, *diskFlipProb, *linkReliable, *maxRetries)

	src := genData(*sizeMB)
	srcHash := sha256.Sum256(src)
	fmt.Printf("[APP] 源数据哈希: %x\n", srcHash)

	ok := false
	for attempt := 1; attempt <= *maxRetries; attempt++ {
		fmt.Printf("\n[APP] 尝试 #%d 进行端到端传输…\n", attempt)
		err := transmit(src, *destPath, *packetSize, *linkFlipProb, *diskFlipProb, *linkReliable)
		if err != nil {
			fmt.Printf("[APP] 传输失败: %v（将重试）\n", err)
			continue
		}

		// 端到端校验：读回磁盘，再次计算哈希
		got, err := os.ReadFile(*destPath)
		if err != nil {
			fmt.Printf("[APP] 读取目标文件失败: %v（将重试）\n", err)
			continue
		}
		gotHash := sha256.Sum256(got)
		fmt.Printf("[APP] 目标数据哈希: %x\n", gotHash)

		if bytes.Equal(srcHash[:], gotHash[:]) {
			fmt.Println("[APP] ✅ 端到端校验通过：传输成功")
			ok = true
			break
		}
		fmt.Println("[APP] ❌ 端到端校验失败：数据不一致（将重试）")
	}

	if !ok {
		fmt.Println("\n[APP] ❌ 多次重试后仍失败。结论：无论链路多可靠，只有端到端校验+重试才能最终保证正确性。")
	} else {
		fmt.Println("\n[APP] ✅ Demo 完成：展示了端到端校验在系统正确性中的决定性作用。")
	}
}

// ---------------------- 核心流程 ----------------------

func transmit(src []byte, dest string, pktSize int, pFlip, pDisk float64, linkReliable bool) error {
	pkts := chunk(src, pktSize)
	buf := make([]byte, 0, len(src))

	for i, p := range pkts {
		if linkReliable {
			// 链路层实现“每包校验+重传”
			wp, tries := "", 0
			for {
				tries++
				cand := maybeFlip(p, pFlip)
				if verifyPacket(cand) {
					buf = append(buf, stripCRC(cand)...)
					wp = fmt.Sprintf("[LINK] 包#%d 成功（重传次数=%d）", i, tries-1)
					break
				}
				if tries%10 == 0 {
					fmt.Printf("[LINK] 包#%d 多次校验失败，继续重试…\n", i)
				}
			}
			if wp != "" && i%200 == 0 {
				fmt.Println(wp)
			}
		} else {
			// 无链路可靠性：可能携带错误直接上交应用
			cand := maybeFlip(addCRC(p), pFlip)
			if verifyPacket(cand) {
				buf = append(buf, stripCRC(cand)...)
			} else {
				// 校验失败仍然“交付”，模拟不可靠链路
				buf = append(buf, stripCRC(cand)...)
			}
		}
	}

	// 写盘
	if err := os.WriteFile(dest, buf, 0644); err != nil {
		return err
	}
	// 磁盘层面腐化（强调：即使链路可靠，存储仍可能损坏）
	if rand.Float64() < pDisk && len(buf) > 0 {
		f, err := os.OpenFile(dest, os.O_RDWR, 0644)
		if err != nil {
			return err
		}
		defer f.Close()
		pos := rand.Intn(len(buf))
		flipByteInFile(f, int64(pos))
		fmt.Printf("[DISK] ⚠️ 写盘后模拟腐化：pos=%d\n", pos)
	}
	return nil
}

// ---------------------- 工具函数 ----------------------

func genData(mb int) []byte {
	total := mb * 1024 * 1024
	var b bytes.Buffer
	for b.Len() < total {
		chunk := []byte("EndToEnd-Argument@MIT: correctness lives at the endpoints.\n")
		b.Write(chunk)
	}
	return b.Bytes()[:total]
}

func chunk(b []byte, n int) [][]byte {
	var out [][]byte
	for i := 0; i < len(b); i += n {
		end := i + n
		if end > len(b) {
			end = len(b)
		}
		out = append(out, b[i:end])
	}
	return out
}

func maybeFlip(p []byte, prob float64) []byte {
	out := make([]byte, len(p))
	copy(out, p)
	if rand.Float64() < prob && len(out) > 0 {
		i := rand.Intn(len(out))
		out[i] ^= 0xFF // 粗暴翻转一个字节
	}
	return out
}

// 简单 CRC：这里用 sha256 的前 4 字节作为“轻量校验”
func addCRC(p []byte) []byte {
	h := sha256.Sum256(p)
	crc := h[:4]
	return append(p, crc...)
}

func verifyPacket(p []byte) bool {
	if len(p) < 4 {
		return false
	}
	data := p[:len(p)-4]
	crc := p[len(p)-4:]
	h := sha256.Sum256(data)
	return bytes.Equal(crc, h[:4])
}

func stripCRC(p []byte) []byte {
	if len(p) < 4 {
		return p
	}
	return p[:len(p)-4]
}

func flipByteInFile(f *os.File, pos int64) {
	_, _ = f.Seek(pos, io.SeekStart)
	var one [1]byte
	if _, err := f.Read(one[:]); err != nil {
		return
	}
	one[0] ^= 0xFF
	_, _ = f.Seek(pos, io.SeekStart)
	_, _ = f.Write(one[:])
}

// （可选）把一些数值编码到日志里，方便外部复现实验
func putU32(b *bytes.Buffer, v uint32) {
	var tmp [4]byte
	binary.BigEndian.PutUint32(tmp[:], v)
	b.Write(tmp[:])
}
