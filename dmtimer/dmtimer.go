package dmtimer

import "time"

type DmTimers map[string]*time.Timer

func (d *DmTimers) Add(key, value string){
	mu.Lock()
 	d[key] = value
	mu.Unlock()	
}

func (d *DmTimers) Del(key string){
	mu.Lock()
	delete(d, key)
	mu.Unlock()	
}


func ParseTimerID(url string) (string, error){
	// expeted format foo.com/ping/123456asdf we want the part afer /ping/
	urlParts := strings.Split(url, "/")
	if len(urlParts) <= 1 {
		return "", fmt.Errorf("Url did not include a Timer ID")
	}
	
	// get rid of any stray spaces
	timerID = string.Trim(timerID)
	
	return timerID, nil
} 