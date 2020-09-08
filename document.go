package main

import (
  "bufio"
  "fmt"
  "log"
  "regexp"
  "strings"
  "strconv"
	"os"
  "path/filepath"
  "github.com/PuerkitoBio/goquery"
)

// document represents a tbc blog post
type document struct {
	Title string
  Text  string
	ID    int
}

func visit(filenames *[]string) filepath.WalkFunc {
  return func(path string, info os.FileInfo, err error) error {
    if err != nil {
      log.Fatal(err)
    }
    if info.IsDir() {
      return nil
    }
    fn := filepath.Base(path)
    first_char := fn[0]
    if first_char != '0' && first_char != '1' {
      return nil
    }
    i := strings.LastIndex(fn, ".")
    if -1 == i {
      return nil
    }
    if i > len(fn) - 4 {
      return nil
    }
    i += 1
    ext := fn[i:i+3]
    if ext != "htm" {
      //fmt.Println(ext)
      return nil
    }
    *filenames = append(*filenames, fn)
    return nil
  }
}

func scanurls(path string) ([]string, error) {

  //pattern := regexp.MustCompile("<tr><td align="right">\d{4}</td><td>\d{4}-\d{2}-\d{2}</td><td><a href="(http[^"]*)">([^\<]*)</a>&nbsp;&nbsp;&nbsp;<a href="([^"]*)">^</a>&nbsp;&nbsp;</td><td>[^\<]*</td></tr>")  

  pattern := regexp.MustCompile("<tr><td align=\"right\">\d{4}</td><td>\d{4}-\d{2}-\d{2}</td><td><a href=\"(http[^\"]*)\">[^\<]*</a>&nbsp;&nbsp;&nbsp;<a href=\"[^\"]*\">^</a>&nbsp;&nbsp;</td><td>[^\<]*</td></tr>")
 
  file, err := os.Open(path)
  if err != nil {
    return nil, err
  }
 
  defer file.Close()
 
  scanner := bufio.NewScanner(file)
 
  scanner.Split(bufio.ScanLines) 
 
  var urls []string
 
  for scanner.Scan() {
    line := scanner.Text()
    url := pattern.FindString(line)
    if( 4 < len(url) {    
      urls = append(urls, url)
  }
  return urls, nil
}

func loadDocuments(path string) ([]document, error) {
  
  // List the blog post source files
  
  var filenames []string
  err := filepath.Walk(path, visit(&filenames))
  if err != nil {
    panic(err)
  }
  n := len(filenames)
  fmt.Println(n, "files")
  
  // Load URLs from index.html
  
  urls, err := scanurls(path + "/index.html")
  if err != nil {
    panic(err)
  }
  for i, url := range urls {
    fmt.Println(i, url)
  }

  // Retrieve blog post document content
  
  docs := make([]document,0,n)

  return docs, nil
  
  for _, fn := range filenames {

    fp := path + "/" + fn

    f, err := os.Open(fp)
    if err != nil {
      log.Fatal(err)
    }
    //defer func() {
    //  if err = f.Close(); err != nil {
    //    log.Fatal(err)
    //  }
    //}()

    r := bufio.NewReader(f)    
    doc, _ := goquery.NewDocumentFromReader(r)
    id, _ := strconv.Atoi(fn[0:4])
    if(!(0<id)){
      log.Fatal("Expected positive blog post number, not ", id)
    }
    title := doc.Find("h3").First().Text()
    docs = append(docs, document{title,doc.Text(),id-1})
    f.Close()
  }
  for i := range docs {
    if( docs[i].ID != i ) {
      fmt.Println(i, docs[i].ID, docs[i].Title)
      log.Fatal("Doc index out of sync with blog post number")
    }
	}
	return docs, nil
}
