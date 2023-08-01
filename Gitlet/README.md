# Gitlet Design Document

**By**: Iman Kahssay

## Overview
Using IntelliJ (and writing in Java), I implemented a version-control system (which is essentially a backup system for related collections of files) that mimics some of the basic features of the popular system called Git. The main functionality that this project, Gitlet, supports are: saving contents of entire directories of files (i.e. <b><i>committing</i></b> files in Git), stopping files from being committed and removing those files from the working directory (i.e. <b><i>rm</i></b> or <b><i>removing</i></b> in Git), restoring a version of one or more files or entire commits (i.e. <b><i>checking out</i></b> files in Git), viewing the history of the user’s backups (i.e. <b><i>logging</i></b> in Git), maintaining related sequences of commits (called <i>committing branches</i> in Git), and merging changes made from one branch into another (i.e. <b><i>merging</i></b> in Git).

## Internal Structures
Real Git distinguishes several different kinds of objects. For our purposes, the important ones are

- <b><i>blobs:</i></b> The saved contents of files. Since Gitlet saves many versions of files, a single file might correspond to multiple blobs: each being tracked in a different commit.
- <b><i>trees:</i></b> Directory structures mapping names to references to blobs and other trees (subdirectories).
- <b><i>commits:</i></b> Combinations of log messages, other metadata (commit date, author, etc.), a reference to a tree, and references to parent commits. The repository also maintains a mapping from branch heads to references to commits, so that certain important commits have symbolic names.
Gitlet simplifies from Git still further by

- Incorporating trees into commits and not dealing with subdirectories (so there will be one “flat” directory of plain files for each repository).
- Limiting ourselves to merges that reference two parents (in real Git, there can be any number of parents.)
- Having our metadata consist only of a timestamp and log message. A commit, therefore, will consist of a log message, timestamp, a mapping of file names to blob references, a parent reference, and (for merges) a second parent reference.

Every object – every blob and every commit in this case – has a unique integer id that serves as a reference to the object. An interesting feature of Git is that these ids are universal: unlike a typical Java implementation, two objects with exactly the same content will have the same id on all systems (i.e. my computer, your computer, and anyone else’s computer will compute this same exact id). In the case of blobs, “same content” means the same file contents. In the case of commits, it means the same metadata, the same mapping of names to references, and the same parent reference. The objects in a repository are thus said to be content addressable.

Both Git and Gitlet accomplish this the same way: by using a cryptographic hash function called SHA-1 (Secure Hash 1), which produces a 160-bit integer hash from any sequence of bytes. Cryptographic hash functions have the property that it is extremely difficult to find two different byte streams with the same hash value (or indeed to find any byte stream given just its hash value), so that essentially, we may assume that the probability that any two objects with different contents have the same SHA-1 hash value is 2-160 or about 10-48. Basically, we simply ignore the possibility of a hashing collision, so that the 
The SHA-1 hash value, rendered as a 40-character hexadecimal string, also makes a convenient file name for storing data in the .gitlet directory. It also gives us a convenient way to compare two files (blobs) to see if they have the same contents: if their SHA-1s are the same, we simply assume the files are the same.

Reading and writing the internal objects from and to files is actually pretty easy, thanks to Java’s serialization facilities. The term “serialization” refers to the conversion from some arbitrary structure (array, tree, graph, etc.) to a serial sequence of bytes. Once an object is serialized, it can be saved in a file *persistly*, meaning the object and its data will be saved and continue to exist even after the system that runs has ended or turned off.

As you can see below, each commit (rectangle) points to some blobs (circles), which contain file contents. The commits contain the file names and references to these blobs, as well as a parent link. These references, depicted as arrows, are represented in the .gitlet directory using their SHA-1 hash values (the small hexadecimal numerals above the commits and below the blobs). The newer commit contains an updated version of wug1.txt, but shares the same version of wug2.txt as the older commit. The Commit class will store all of the information that this diagram shows.

<img src="images/commits and blobs diagram.png" alt="Commits and Blobs Diagram" align="middle" width="350px" >


## Classes and Data Structures

### Class 1: Main

This is the entry point to the program.
It takes in arguments from the command line and based on the 
command (the first element of the args array) calls the 
corresponding command in Repository which will 
actually execute the logic of the command. 
It also validates the arguments based on the command to 
ensure that enough arguments were passed in.

#### Fields

This class has no fields and hence no associated state: it 
simply validates arguments and defers the execution to the 
Repository class.


### Class 2: Repository

This is where the main logic of the program will live. This 
class will handle all of the actual Gitlet commands by 
reading/writing from/to the correct file, setting up 
persistence, and additional error checking.

It will also be responsible for setting up all persistence 
within Gitlet. This includes creating the .gitlet folder as 
well as the folder and file where all the Commit objects and Blobs are stored.

This class defers all Commit specific logic to the Commit class: 
for example, instead of having the Repository class handle 
Commit serialization and deserialization, the Commit 
class will contain the logic for that.


#### Important Fields

1. Field 1: **static final File CWD = new File(System.getProperty("user.dir"))**
   
    The Current Working Directory. Since it has the 
    package-private access modifier (i.e. no access modifier), 
    other classes in the package may use this field. 
    This is useful for the other File objects used in this project.

2. Field 2: **static final File GITLET_REPOSITORY_FOLDER = Utils.join(CWD, ".gitlet")** 
   
    The hidden .gitlet directory. This is where all of the 
    states of the Repository will be stored, including additional 
    things like the Commit objects and the Blobs. It is also 
    package private as other classes will use it to store their state.

3. Field 3: **public File HEAD = join(Repository.GITLET_REPOSITORY, "HEAD")**

   'Points' to the most current commit (the 'leaf') by saving the
    SHA-1 Hash of that Commit (i.e. commitID) in a text file.
   
4. Field 4: **public static final File MASTER = join(Repository.GITLET_REPOSITORY, "MASTER")**

   'Points' to the 'root' commit by saving the SHA-1 Hash of 
   that commit (i.e. commitID) (therefore, saving the branch/tree 
   of the repository).

These fields are both static since the Repository class does not actually instantiate a Repository object: it is simply used to house functions. 
If there were additional non-static states (like the Dog class), it would need to be serialized and saved in to a file.



### Class 3: Commit

This class represents a Commit that will be stored in a file. 
Because each Commit will have a unique ID, that unique ID can 
simply be used as the name of the file that the object is serialized to.


All Commit objects are serialized within the COMMITS_FOLDER which 
is within the GITLET_REPOSITORY_FOLDER. The Commit class has 
helpful methods that will return the Commit object corresponding 
to some String ID given to it, as well as write that Commit to 
a file to persist its changes.

#### Important Fields

1. Field 1: **static final File COMMITS_FOLDER = Utils.join(Repository.GITLET_REPOSITORY_FOLDER, ".commits")** 
   
   The File object that corresponds to the directory containing 
   all the serialized Commit objects. This is static since all 
   Commit objects are stored within the same directory. 

2. Field 2: **private long timestamp** -> The time Commit object was "committed" (or created).

3. Field 3: **private String message** -> User's commit message when creating a Commit.

4. Field 4: **private String commitID** -> The SHA1 hashvalue of Commit (will also be the name of the File Object that will save the serialized Commit's Object's data).

5. Field 5: **private String commitParentID** -> The SHA1 hashvalue of previous Commit (i.e. commit's parent).

6. Field 6: **private Hashtable<String, String> tableOfBlobs;** ->  Hashtable of Blob references (the SHA1 value of Blobs (it's ID)).

### Class 4: Blob

This class represents a Blob that will be stored in a file.
Because each Blob will have a unique ID, that unique ID may simply be used 
as the name of the file that the object is serialized to.
This class will save the different *versions* of a single file.

#### Important Fields

1. Field 1 : **static final File BLOBS_FOLDER = Utils.join(Repository.GITLET_REPOSITORY, ".blobs")**
   
    The File object that corresponds to the directory 
    containing all the serialized Blob objects.
    This is static since all Blob objects are stored within 
    the same directory.
   
2. Field 2: **private String blobID** -> The SHA1 value of this Blob (it's ID).

3. Field 3: **private String fileName** -> Name of the file Blob will be saving the contents of.
  
4. Field 4: **private File fileContents** -> Saves the target file's version/contents into a text file.
    - Thought about saving it as a String, but the contents can be too large.

### Class 5: Utils

This class contains helpful utility methods to read/write 
objects or String contents from/to files, as well as reporting 
errors when they occur.

### Class 6: Staging Area

This class represents the Staging Area that stores files before
they are committed (in *add* and *remove* commands).

#### Important Fields

1. Field 1 :


### Class 8: Gitlet Exception

This class represents all of the exceptions and/or error 
messages that may be thrown.


## Algorithms

## Persistence

The directory structure looks like this:

    CWD                                         <==== Whatever the current working directory is.
    └── .gitlet                                 <==== All persistant data is stored within here
        ├── commits                             <==== All commits are stored in this directory
            ├── commit1_ID                      <==== A single Commit instance stored to a file
            ├── commit2_ID
            ├── ...
            └── commit_n_ID
        └── blobs                               <==== All blobs are stored in this directory
            ├── blob1_ID                        <==== A single Blob instance stored to a file
                ├── fileContents for Blob1_ID       <==== The fileContents saved in a file in a Blob object.
            ├── blob2_ID 
                ├── fileContents for Blob2_ID 
            ├── ...
            ├── blob_n_ID
                └── fileContents for Blob_n_ID 
        └── staging area directory                <==== Persistent Staging Area class is stored in this
            ├── state.txt                         <==== The current state of a commit is saved in a text file
        └── HEAD.txt                              <==== Persistent Head Pointer saved in a text file.
       

