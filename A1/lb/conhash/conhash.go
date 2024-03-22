package conhash

import (
<<<<<<< HEAD
	"fmt"
)

=======
	// "crypto/sha256" 
	// "encoding/binary"
	"fmt"
)

// func last9BitsSHA256(number uint64) int {
// 	// Convert the number to bytes
// 	dataBytes := make([]byte, 8)
// 	binary.BigEndian.PutUint64(dataBytes, number)

// 	// Calculate the SHA-256 hash
// 	hash := sha256.Sum256(dataBytes)

// 	// Extract the last 9 bits from the hash
// 	last9Bits := int(binary.BigEndian.Uint16(hash[len(hash)-2:]) & 0b111111111)

// 	return last9Bits
// }

>>>>>>> d0baf8e1f9b3f8c7be44e9844733f395d537336c

// Node represents a node in the ConHash structure.
type Node struct {
	Occ  bool
	Name string
}

// ConHash is a consistent hashing implementation.
type ConHash struct {
	HashD      []Node
	Size		int
	VirtServ   	int
	Nserv       int
	AllServers map[string]int
<<<<<<< HEAD
=======
	ServerID map[string]int
>>>>>>> d0baf8e1f9b3f8c7be44e9844733f395d537336c
}

// NewConHash creates a new ConHash instance.
func NewConHash(m, k int) *ConHash {
	return &ConHash{
		HashD:      make([]Node, m),
		Size:       m,
		VirtServ:   k,
		Nserv:		0,
		AllServers: make(map[string]int),
<<<<<<< HEAD
=======
		ServerID: 	make(map[string]int),
>>>>>>> d0baf8e1f9b3f8c7be44e9844733f395d537336c
	}
}

func (c *ConHash) getServHash(i, j int) int {
<<<<<<< HEAD
	val := 1*i*i + 1*j*j + 2*j + 25
=======
	val := i*i + j*j + 2*j + 25
>>>>>>> d0baf8e1f9b3f8c7be44e9844733f395d537336c
	val %= int(c.Size)
	return int(val)
}

func (c *ConHash) getCliHash(i int) int {
<<<<<<< HEAD
	val := 1*i*i + 1*2*i + 17
=======
	val := i*i + 2*i + 17
>>>>>>> d0baf8e1f9b3f8c7be44e9844733f395d537336c
	val %= int(c.Size)
	return int(val)
}

// Add adds servers to the consistent hash ring.
func (c *ConHash) Add(ids []int, Names []string) int {
	if len(ids) != len(Names) {
		return 0
	}

	for _, name := range Names {
		if _, ok := c.AllServers[name]; ok {
			return 0
		}
	}
	
	if (c.Nserv+len(ids))*c.VirtServ >= c.Size {
		return 0
	}
	c.Nserv += len(ids)

	for i := 0; i < len(ids); i++ {
		c.AllServers[Names[i]] = 1
<<<<<<< HEAD
=======
		c.ServerID[Names[i]] = ids[i]
>>>>>>> d0baf8e1f9b3f8c7be44e9844733f395d537336c
		for j := 0; j < c.VirtServ; j++ {
			hash := c.getServHash(ids[i], j)
			for c.HashD[hash].Occ {
				hash = (hash + 1) % c.Size
			}
			c.HashD[hash] = Node{true, Names[i]}
		}
	}
	return 1
}

// GetConfig prints the configuration of the consistent hash ring.
func (c *ConHash) GetConfig() {
	for i := 0; i < c.Size; i++ {
		fmt.Printf("Index: %d Status: %t Server: %s\n", i, c.HashD[i].Occ, c.HashD[i].Name)
	}
}

// AddServer adds a single server to the consistent hash ring.
func (c *ConHash) AddServer(id int, Name string) int {
	if _, ok := c.AllServers[Name]; ok {
		fmt.Println("Same server name already exists")
		return 0
	}

	if (c.Nserv+1)*c.VirtServ >= c.Size {
		fmt.Println("Size limit exceeded")
		return 0
	}

	c.Nserv++
	c.AllServers[Name] = 1
<<<<<<< HEAD
=======
	c.ServerID[Name] = id
>>>>>>> d0baf8e1f9b3f8c7be44e9844733f395d537336c

	for j := 0; j < c.VirtServ; j++ {
		hash := c.getServHash(id, j)
		for c.HashD[hash].Occ {
			hash = (hash + 1) % c.Size
		}
		c.HashD[hash] = Node{true, Name}
	}

	return 1
}

// RemoveServer removes a server from the consistent hash ring.
func (c *ConHash) RemoveServer(Name string) int {
	if _, ok := c.AllServers[Name]; !ok {
		return 0
	}

<<<<<<< HEAD
	for i := range c.HashD {
		if c.HashD[i].Name == Name {
			c.HashD[i] = Node{false, ""}
		}
	}
	delete(c.AllServers, Name)
=======
	for j := 0; j < c.VirtServ; j++ {
		hash := c.getServHash(c.ServerID[Name], j)
		for c.HashD[hash].Name != Name {
			hash = (hash + 1)%c.Size
		}
		c.HashD[hash] = Node{false, ""}
	}

	delete(c.AllServers, Name)
	delete(c.ServerID, Name)
>>>>>>> d0baf8e1f9b3f8c7be44e9844733f395d537336c
	c.Nserv--
	return 1
}

// GetServer returns the server for the given client ID.
func (c *ConHash) GetServer(id int) string {
	if c.Nserv == 0 {
		return "No Server Allocable"
	}
	hash := c.getCliHash(id)
<<<<<<< HEAD
	fmt.Println("hash for ", id, " --> ", hash)
=======
	// fmt.Println("hash for ", id, " --> ", hash)
>>>>>>> d0baf8e1f9b3f8c7be44e9844733f395d537336c
	hash = (hash + 1) % c.Size
	for !c.HashD[hash].Occ {
		hash = (hash + 1) % c.Size
	}

	return c.HashD[hash].Name
}
