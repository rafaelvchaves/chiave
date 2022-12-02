scp -i ../../../chiave-ssh-key.pem ec2-user@ec2-3-231-214-28.compute-1.amazonaws.com:chiave/replica/leader/cpu.pprof .
pprof -http :8080 cpu.pprof
