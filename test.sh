#! /bin/bash

function main() {

  for i in {1..10000}; do 
    cat hello | openssl s_client -connect localhost:8090
  done 

}

main
