package main

import (
	"encoding/gob"
	"fmt"
	"os"

	//"encoding/gob"
	"github.com/spaolacci/murmur3"
	"hash"
	"math"
	//"os"
	"time"
)

type BloomFilter struct {
	M     uint          // Velicina Set-a
	K     uint          // Broj hash funkcija
	P     float64       // False-positive vjerovatnoca
	Set   []byte        // Set sa bitovima
	hashs []hash.Hash32 // hash funkcije
	TimeConst uint
}

func CreateBloomFilter(n uint, p float64) *BloomFilter {
	m := CalculateM(int(n), p)
	k := CalculateK(int(n), m)
	hashs, tc := CreateHashFunctions(k)
	bf := BloomFilter{m, k, p, make([]byte, m), hashs, tc}
	return &bf
}

func (bf *BloomFilter) Add(elem Element) {
	for _, hashF := range bf.hashs {
	i := HashIt(hashF, elem.key, bf.M)
	bf.Set[i] = 1
	}
}

func (bf *BloomFilter) Query(elem string) bool {
	for _, hashF := range bf.hashs {
		i := HashIt(hashF, elem, bf.M)
		if bf.Set[i] != 1 {
			return false
		}
	}
	return true
}

func HashIt(hashF hash.Hash32, elem string, m uint) uint32 {
	_, err := hashF.Write([]byte(elem))
	if err != nil {
		panic(err)
	}
	i := hashF.Sum32() % uint32(m)
	hashF.Reset()
	return i
}

func CalculateM(expectedElements int, falsePositiveRate float64) uint {
	return uint(math.Ceil(float64(expectedElements) * math.Abs(math.Log(falsePositiveRate)) / math.Pow(math.Log(2), float64(2))))
}

func CalculateK(expectedElements int, m uint) uint {
	return uint(math.Ceil((float64(m) / float64(expectedElements)) * math.Log(2)))
}

func CreateHashFunctions(k uint) ([]hash.Hash32, uint) {
	var h []hash.Hash32
	ts := uint(time.Now().Unix())
	for i := uint(0); i < k; i++ {
		h = append(h, murmur3.New32WithSeed(uint32(ts+1)))
	}
	return h, ts
}

func CopyHashFunctions(k uint, tc uint) ([]hash.Hash32) {
	var h []hash.Hash32
	for i := uint(0); i < k; i++ {
		h = append(h, murmur3.New32WithSeed(uint32(tc+1)))
	}
	return h
}

func writeBloomFilter(filename string, bf *BloomFilter) {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(bf)
	if err != nil {
		panic(err)
	}
}

func readBloomFilter(filename string) (bf *BloomFilter) {
	file, err := os.Open(filename)
	if err != nil {
		return nil
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	bf = new(BloomFilter)
	_, err = file.Seek(0, 0)
	if err != nil {
		return nil
	}

	for {
		err = decoder.Decode(bf)
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println(*bf)
	}
	bf.hashs = CopyHashFunctions(bf.K, bf.TimeConst)
	return
}

//func main() {
//
//	bf := CreateBloomFilter(30, 2)
//	bf.Add("Bojan")
//	bf.Add("Mićo")
//	bf.Add("Katarina")
//	bf.Add("Milica")
//	bf.Add("Miloš")
//	fmt.Println("\nNemanja ? ", bf.Query("Nemanja"))
//	fmt.Println("Katarina ? ", bf.Query("Katarina"))
//	bf.Add("Branko")
//	bf.Add("Gaga")
//	bf.Add("Djuro")
//	bf.Add("Suncica")
//	bf.Add("Ljupka")
//	bf.Add("Krinka")
//	bf.Add("Djole")
//	bf.Add("Mirjana")
//	bf.Add("Jovo")
//	bf.Add("Dado")
//	bf.Add("Mira")
//	fmt.Println("\nNemanja ? ", bf.Query("Nemanja"))
//	fmt.Println("Jovo ? ", bf.Query("Katarina"))
//
//	fmt.Println("\nSerialization in progress ...")
//
//	nwf, _ := os.Create("bf.gob")
//	nwf.Close()
//
//	file, _ := os.OpenFile("bf.gob", os.O_RDWR, 0666)
//	defer file.Close()
//	encoder := gob.NewEncoder(file)
//	err := encoder.Encode(bf)
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	decoder := gob.NewDecoder(file)
//	var srs = new(BloomFilter)
//	file.Seek(0, 0)
//	for {
//		err = decoder.Decode(srs)
//		if err != nil {
//			fmt.Println(err)
//			break
//		}
//		fmt.Println(*srs)
//	}
//	srs.hashs = CopyHashFunctions(srs.K, srs.TimeConst)
//	fmt.Println("\nNemanja ? ", srs.Query("Nemanja"))
//	fmt.Println("Jovo ? ", srs.Query("Katarina"))
//}
