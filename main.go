// BitBackup - Backup & Sync your bitbucket team repositories
// Balamurali Pandranki - <balamurali[at]live.com>
package main

import (
	"flag"
	"os/exec"
	"sync"

	"github.com/ktrysmt/go-bitbucket"
	log "github.com/sirupsen/logrus"
)

func main() {
	username := flag.String("username", "", "Bitbucket Username")
	password := flag.String("password", "", "Bitbucket Password")
	team := flag.String("team", "utuindia", "Bitbucket Team name")
	backupdir := flag.String("backup-dir", "bitbucket-backup", "Folder to backup")
	flag.Parse()
	println(`// BitBackup - Backup & Sync your bitbucket team repositories
// Balamurali Pandranki - <balamurali[at]live.com>`)
	_, err := exec.LookPath("git")
	if err != nil {
		log.Fatalln("Git is not available on the system, Please install Git from https://git-scm.com/")
	}

	if *username == "" && *password == "" {
		log.Fatalln("Invalid Username/Password")
	}
	if *backupdir == "" {
		log.Fatalln("Invalid Backup directory")
	}
	c := bitbucket.NewBasicAuth(*username, *password)
	//Get all bitbucket repositories
	res, err := c.Teams.Repositories(*team)
	if err != nil {
		log.Fatalln("Bitbucket Error: ", err)
	} else {
		resp := res.(map[string]interface{}) //receive the data as json format
		if err != nil {
			log.Errorln(err)
		}
		vals := resp["values"]
		values := vals.([]interface{})
		log.Infoln("Found", len(values), "Respositories to Sync")
		//Async
		wg := sync.WaitGroup{}
		for _, val := range values {
			repo := val.(map[string]interface{})
			//links := repo["links"].(map[string]interface{})
			//repolink := links["clone"].([]interface{})[0].(map[string]interface{})["href"].(string)
			wg.Add(1)
			go func() {
				log.Infoln("Cloning Repo", repo["name"], "into", repo["slug"])
				giturl := "https://" + *username + ":" + *password + "@bitbucket.org/" + *team + "/" + repo["slug"].(string) + ".git"
				println("Git Clone URL:", giturl)
				cmd := exec.Command("git", "clone", giturl, *backupdir+"/"+repo["slug"].(string))
				err := cmd.Run()
				if err != nil {
					//Something went wrong
					log.Errorln(err)
				} else {
					log.Infoln("Done Cloning Repo", repo["name"], "into", repo["slug"], "directory")
				}
				wg.Done()
			}()
		}
		wg.Wait()
	}
}
