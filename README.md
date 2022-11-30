This programs demonstrate the filestore capabilites as:
Adding the files to the store (while also checking the duplication)
Deleting the file
Updating the file
Listing all the files
Counting total word count
Showing n most-frequent or less-frequent words in the files

Example commands:
filestore add file1.txt file2.txt
filestore ls
filestore rm file1.txt
filestore update file1.txt
filestore wc
filestore freq-words -n 10 --order=asc


The app has been dockerized and pushed to docker hub as "docker.io/ankitrkharb/filestore:1.0"
Same docker image has been pulled from docker hub and checked by deploying to local kubernetes cluster.
