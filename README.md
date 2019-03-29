## quick start
###  interactive

simplessh


        type h to get help info

        >>> h
        
        c id/alias to connect a cfg
        
        n to new a cfg
        
        nc to new a cfg and connect it
        
        u id/alias to update a cfg
        
        uc id/alias to update a cfg and connect it
        
        l to list cfg
        
        d id/alias to del a cfg
        
        q to quit
        
        h to get help info

        >>>


### one time
#### 新建

simplessh n

        addr:(192.168.1.100:22) 192.168.1.100:22


        user:(default root) root


        password: test

        alias: test

#### 新建并链接
simplessh nc

        addr:(192.168.1.100:22) 192.168.1.100:22


        user:(default root) root


        password: test

        alias: test

#### 更新

通过序号更新

simplessh u 2

        addr:(192.168.1.100:22) 192.168.1.100:22


        user:(default root) root


        password: test

        alias: test


通过别名更新

simplessh u test

        addr:(192.168.1.100:22) 192.168.1.100:22


        user:(default root) root


        password: test

        alias: test


#### 更新并链接

通过序号更新并链接

simplessh uc 2

        addr:(192.168.1.100:22) 192.168.1.100:22


        user:(default root) root


        password: test

        alias: test


通过别名更新并链接

simplessh uc test

        addr:(192.168.1.100:22) 192.168.1.100:22


        user:(default root) root


        password: test

        alias: test

#### 显示

simplessh l


        0)      192.168.1.100:22

#### 链接

通过序号链接

simplessh c 0

通过别名链接

simplessh c test

####  删除

通过序号删除

simplessh d 0

通过别名删除

simplessh d test


#### 显示具体配置

通过序号显示

simplessh i 0

通过别名显示

simplessh i test