package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
)

func (c *ConfigStruct) NotifyConfigFile() error {
	if c.watcher != nil {
		return nil
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	// Add a path.
	err = watcher.Add(c.configPath)
	if err != nil {
		return err
	}

	var stop = make(chan error)

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					stop <- nil
					return
				}
				if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {
					err := c.Reload()
					if err != nil {
						fmt.Printf("Config file reload error: %s\n", err.Error())
					} else {
						fmt.Printf("Config file reload success\n")
					}
				} else if event.Has(fsnotify.Remove) || event.Has(fsnotify.Rename) {
					fmt.Printf("Config file has been remove\n")
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					stop <- nil
					return
				}
				fmt.Printf("Config file notify error: %s\n", err.Error())
			}
		}
	}()

	return nil
}

func (c *ConfigStruct) CloseNotifyConfigFile() {
	if c.watcher == nil {
		return
	}

	_ = c.watcher.Close()
}
