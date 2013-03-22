# Tin Foil Hat

TFH is a file storage system that builds on top of a modified version of the Kademlia distributed hash table (http://xlattice.sourceforge.net/components/protocol/kademlia/specs.html). It takes a file as an input, and encrypts that file. To decrypt, the user inputs the file path of the key and the filepath of the place to store the unencrypted file. TFH gives you the ability to securely store your files on a network that is distributed, fault tolerant, and offers better safety* than standard AES encryption.

WARNING! This program is an EARLY EARLY ALPHA stage, and probably does not offer any sort of security at all in its current state!! DO NOT use this in any sort of production enviroment. 


## Way way way safer than plain AES
TFH makes it completely intractable to brute force decrypt the original file. As an example, if a malicious agent were to detect just 20 chunks (or 640 bytes) coming from your machine, it would know that <=10% (rounded up) of these are junk so the attacker would know that 18, 19 or all 20 of these chunks are real. In that case it would then have to check every subset of size 18, 19 or 20 and then every permutation of each of those subsets just to get to a point where AES decryption can be attempted. Even if the malicious agent could solve AES decryption in a tractable amount of time, they would have to make 20! + (20 Choose 19) *19! + (20 Choose 18) *18! ~= 6.0*10^18  different attempts to crack every possible combination so expected 3 * 10^18. So even if we make the absurd assumption that the attacker could crack an aes encrypted file in one second, it would still take 9.6*10^10 years or approximately 7 times the estimated age of the universe to crack. Remember this is just for ~700 bytes of data (it takes even longer as the file size grows)!

## Building and Running
- `go install tin-foil-hat/tfh` to install the binary.
- Then, you can run it via `./tfh ip_and_port_to_bind_to ip_and_port_of_existing_node_to_connect_to`
- This will boot up one node. `./tfh 127.0.0.1:8090 127.0.0.1:8090` is an example for the first node.

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
- Our system will interleave a random amount up to 10% “Junk Chunks” to confuse anyone who might be trying to reconstruct the file. Only the transmitter of the file knows which segments are real and which ones are junk

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

## TODO

We *think* this system is pretty nifty, but want to talk to people who are experts in Go and in Crypto to let us know what they think! If there is enough interest in the concept, we'd love to make the system robust, and production ready.

- TLS
- Node Reputation system (so that it can be used with untrusted peers)
- TESTING!!
