package main

import (
  "flag"
)

func main() {
  word := flag.String("word", "", "A word passed to SAOL.")
  flag.Parse()

  if (*word == "") {
    panic("No word provided.")
  }
  // TODO: call the SAOL module with the word
}
