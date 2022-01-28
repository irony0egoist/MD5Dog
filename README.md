# MD5Dog

MD5Dog - Use Golang to Crack Hash

```
    
    Options:
      -H, --hash <data>             Hash value
      -s, --salt <data>             Salt value
      -p, --path <data>             Path to passwords dictionary
      -t, --type <data>             Hash type
      -c, --concurrency <data>      Concurrency count AND file buffer size
    
    Usage:
    *md5(password.salt):
            Example: MD5Dog.exe -t 10 -c 1000 -s abc -H 1b7ff998949c08bfa0d399d41aa0
    cbdf -p dic\pass.txt
    *hmac.md5(password) ;secret=md5(password):
            Example: MD5Dog.exe -t 20 -c 1000 -H 92d858ce796c86d55090ce1f1bb7be9a -p dic\pass.txt

```

Successful Output:
you can see successful line
```
2022/01/28 17:32:35 load passwords: 487 lines
2022/01/28 17:32:35 successful: admin123
2022/01/28 17:32:35 md5Dog finished! 
```
