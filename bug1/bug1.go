package bug1
import "sync" 
// Using mutex to synchronize

// Counter stores a count.
type Counter struct {
	n int64
	mu sync.Mutex
}

// Inc increments the count in the Counter.
func (c *Counter) Inc() {
	c.mu.Lock()
	c.n++
	c.mu.Unlock()
}
