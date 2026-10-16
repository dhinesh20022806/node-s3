### Node s3

This software lives in the node and recive or send objects to [data routine service](https://github.com/dhinesh20022806/data_routing_service)

## Write object

# Design choice

store each object on single file leads to inode exhausted and large memory files to load up in to memory can leads to crash the program. so i choose in gb size files it's 4 Gb

due to every request needs to write on single file it's imporatant to make the data serilaizable , make sure concurrent request overwrite write object on the files it's restrict the write throughput, to fix those issues didcated files owned by each go routine, go routine for writing files up to no of physical cores present in system. each go routine write their files. let's me break each component in write object.

I will explain top down apporach

dispatcher - this is the responsible for receive write request from data routine service and assign jobs to worker in round robin algorithm.

worker - this is main component that write on file, workers are fired up even before request comes in and owns their files based on incremenatal number initial it's zero and it's presitience stored on disk when program about to crash. these increamental atomic variable so it's guarantee that each worker owns different files . it's check file before write the object to file. if file exceed the 4GB or it's maked as read only files and create new one then write the file. each worker have upper bound to prevent oom error.

after write the file, the worker return Result in go channels about any err, filename - where it's stored, starting address , length of bytes. then store the meta information on sqlite

fileAllocator - which is responsbile for allocate the files to worker without same file access by two write worker by enable mutx locking

filenamegenerator - which is responsible generating name files in serial way
like a-z, a1-z1, .... an-zn.

## Read Request

it's simple compare the write request lookups the sqlite against object key and pick up the file read at address to length of bytes return it.
