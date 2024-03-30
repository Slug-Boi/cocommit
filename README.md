# cocommit

## IMPORTANT NOTICE!!!! 
### Windows defender recently got an update that apparently sees Go Binaries as a trojan... So I guess you have to manually go in and flag it as not being that until Microsoft decides to fix this.
[Link to Github issue talking about this](https://github.com/golang/vscode-go/issues/3182)
Currently built in:
- [x] Go
- [ ] Rust


Currently tested in:
- [x] Go
- [ ] Rust


# How to run
Download the binary for your OS on the [release page](https://github.com/Slug-Boi/cocommit/releases)  
Once downloaded you need to create an alias for your shell guides found below:  
[Bash Guide for alias](https://linuxize.com/post/how-to-create-bash-aliases/)  
*^please note if you use another shell than bash you will need to look up how to do it in that shell^*  
[Powershell Guide for alias](https://stackoverflow.com/questions/24914589/how-to-create-permanent-powershell-aliases)  

Once you've created an alias for the program you need to set an env variable (author_file) in your shell. This should be the path pointing to your author.txt file  
For bash you just need to add this to your .bashrc file:
```
export author_file='path/to/your/aurhor.txt'
```
[Powershell guide for env variable](https://stackoverflow.com/a/714918)

*^Please note that the syntax line at the top of the author.txt file should not be deleted and that you must add at least one author to the file to run the program^*

you can now run it using the alias shorthand you assigned to it 
## Usage:
```
$ cocommit "message" <name1> [name2] [name3]...

$ cocommit "message" <name:email1> [name:email2] [name:email3]...

$ cocommit "message" <name:email1> <name1> [name:email2]

$ cocommit "message" all
*adds all comitters execpt ones tagged with ex*

$ cocommit "message" ^<name1> ^[name2]
*adds all users except the negated users and users tagged with ex*

$ cocommit "message" <group_name>
*adds all users that has that group tag in author file*

$ cocommit
*prints usage*

$ cocommit users
*prints list of users*
```

# Syntax for the author file
The syntax for the author file can be found at the top of the template file included in the repo. It should look like this (opt) is optional syntax:  
```
name_short|Name|Username|email (opt: |ex) (opt: ;;group1 or ;;group1|group2|group3...)
```
opt explained:  
ex -> excludes the given author for all and negation commands  
group -> groups an author which can then be called as an argument to add all people from that group. An author can be a part of multiple groups 

# Why?
Writing co-authors onto commits can be pretty tedious so automating this process as a simple shell alias is a lot nicer

# Workflows
This repo is sort of a test bed for working with Dagger CI but therefore also should support automated testing and building at some point 
*See ci folder for the current dagger workflows*

