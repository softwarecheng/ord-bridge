# ord bridge 
A ordinals bridge written in Go

## Requirements
- Bitcoind node 25
- ord 0.16.0-ordx (https://github.com/softwarecheng/ord/tree/0.16.0-ordx)

### Minimum System Requirements:
- 16+ GB RAM
- 8 Core CPU
- 50+ GB SSD for the indexer db

## Develop
1. Bitcoind node
   ```shell
   `
   `
   ```
2. customized version ord base on 0.16.0
    ```shell
    `
    git clone https://github.com/softwarecheng/ord.git && cd ord && git checkout 0.16.0-ordx && cargo build release
    
    ord --chain testnet --bitcoin-rpc-url 192.168.1.101:18332 \ 
    --data-dir /data2/ord-data/0.16.0/testnet \ 
    --bitcoin-rpc-username jacky --bitcoin-rpc-password 123456 \ 
    --first-inscription-height 2413343 --index-sats \
    ordx export --filename /data/ordx-data-backup/ord-latest/testnet-all-inscription-data1.ordx --cache 5 --update-indexd
    
    `
    ```
   
3. Rename `example.yaml` to `config.yaml` and substitute the values for your system
   go build


