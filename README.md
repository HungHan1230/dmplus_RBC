## 前言
這隻程式是用來做RBC-data framework的簡單模擬。是在論文告急的情況下趕出來的版本，所以非常欠缺程式優化，需要比較大的disk空間來儲存產生的data，運行時間也會很長。

## 運行概念
為了運行RBC-data framework的模擬，大致上的概念是以Random的方式指定某些節點保存blk.dat檔。接著透過一個節點消失加入的pattern來去決定節點的上下線，最後是在某些時間點下去發出request看取得檔案的失敗率（fault rate）是多少。

實際上的細部實作概念如下：

從[Bitnodes](https://bitnodes.io/api)這個網站上下載歷史的Bitcoin節點上線紀錄，在我的實驗中總共下載了8200個network snapshot，期間是4/30 - 5/31

需要計算的部份有三點
- 8200個network snapshot中的churn數量
- 指派哪些節點根據replication factor來儲存blk.dat
- 隨著下載的這個pattern，節點隨著這個pattern加入和離開這個網路，並透過Bitcoin full node的access紀錄來對每個timestamp發起access請求，最後計算失敗率（fault rate）



在Bitnodes中，使用以下Web Api可以得到每個timestamp的節點數
```
https://bitnodes.io/api/v1/snapshots/?page=1000
```

使用另外這個Web Api可以得到每個timestamp下的細部資料
```
https://bitnodes.io/api/v1/snapshots/1588156832
```


## 使用方式

Go語言版本
```
$ go version go1.14.2 linux/amd64
```

主程式是Run.go這隻程式，運行用以下的指令來執行程式
```
$ go run Run.go
```

在%------------------------------Run------------------------------%這個part

建議多開terminal同時去運行

如以下例子，開6個terminal，第一個先執行Baseline_norepair，隨後即註解第1行換第2行執行，以此類推。
```
--------get 500/10G results--------
	RunProcess("Baseline_norepair_",500, false, 10, 1,"10Baseline.csv")
	// RunProcess("R2_",500, false, 10, 2,"10R2.csv")
	// RunProcess("R4_",500, false, 10, 4,"10R4.csv")
	// RunProcess("R8_", 500, false, 10, 8, "10R8.csv")
	// RunProcess("R16_", 500, false, 10, 16, "5R16.csv")
	// RunProcess("R32_", 500, false, 10, 32, "5R32.csv")
```

p.s. Run.go裡的funtion都有寫上註解，但要特別注意路徑設定的部份是在WhoIsAlwaysUp()、CalculateEachSession() 和 RunProcess()中的路徑是原本我電腦的路徑，需要修改。
