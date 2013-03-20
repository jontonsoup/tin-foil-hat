# Tin Foil Hat

TFH is a system that builds on top of a modified version of the Kademlia distributed hash table. It takes a file as an input, and encrypts that file. To decrypt, the user inputs the file path of the key and the filepath of the place to store the unencrypted file. TFH gives you the ability to securely store your files on a network that is distributed, fault tolerant, and offers better safety than standard AES encryption.


## Building and Running
- `go install kademlia-secure/tfh` to install the binary.
- Then, you can run it via `./tfh 127.0.0.1:8090 127.0.0.1:8090`

## Encrypting a file
`encrypt /path/to/file/you/want/to/store /path/to/keyfile/you/want/to/generate`

## Decryping a file
`decrypt /path/to/keyfile /path/to/file/you/want/to/decrypt`


## Advantages of the TFH system:

- AES 256 Encryption of the file
- 32 byte segments
- Distributed storage of segments across multiple physical servers
- Redundant storage of files across multiple physical servers
- Consistency checking of file segments and file as a whole
- Random transmission of file segments so people sniffing the wire cannot reconstruct the file
- This is because both the file segments overlap (meaning there are many segments of a file that repeat)
- And the the order of the file segments used to reconstruct the file is known only to the sender of the file
- Our system will interleave 10% “Junk Chunks” to confuse anyone who might be trying to reconstruct the file. Only the transmitter of the file knows which segments are real and which ones are junk

##Algorithm:
### Encryption

- Input file and path to keyfile
- Randomly generate 32 byte key
- Hash file and save it for later
- Padd bytes so that encrypted file is a multiple of 32 bytes
- Encrypt file with AES 256 encryption
Break file into 32 byte chunks
- Hash each chunk
- Randomly store each chunk in kademlia DHT with key as the hash of the chunk and value as the chunk itself. Also store interleave the storage of 10% as many junk chunks.
- When the remote server receives the chunk it recomputes the hash of the file and checks to make sure that it still matches the hash in the kademlia key.
- Write the bytes padded, the hash of the unencrypted file, and each of the hashes of the chunks in the order that they were created (not in the random order they were sent on the wire) to disk

### Decryption

- Input path to keyfile and path to file you want to generate
- Deserialize the keyfile into bytes padded, the hash of the unencrypted file, and each of the hashes of the chunks in the order that they were created
- Randomly lookup each file chunk in the DHT based on the key(hash) and the junk chunks.
- Reconstruct the original file order
- Trim the padding off the original file
- Check the original file against the hash that was created during the encryption phase

### TODO
- TLS
- Node Repuatation system (so that it can be used with untrusted peers)