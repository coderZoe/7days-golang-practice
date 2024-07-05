在本章中，对于中间件的串行执行，我们使用了一个`Next`方法，这里有必要对这个方法进行详细解释：
```go
func (c *Context) Next() {
	c.index++
	s := len(c.middleware)
	for ; c.index < s; c.index++ {
		c.middleware[c.index](c)
	}
}
```
如上，如果我们假设中间件的写法都是：
```go
func(ctx *gee.Context) {
    //do something
    ctx.Next()
    //do something
}
```
也即每个中间件都调用`ctx.Next()`，那我们其实可以将上述`Next`改为
```go
func (c *Context) Next() {
	c.index++
	if(c.index < len(c.middleware)){
		c.middleware[c.index](c)
	}
}
```
如上不需要循环，因为每个middleWare里都调用了`ctx.Next()`，这个链式调用会一直走下去
当走到最后一个节点，又会由于函数调用入栈的原因，调用结束后会出栈反向执行

但很明显，我们得考虑中间件中无`ctx.Next()`的情况，如果中间件中无`ctx.Next()`，上述版本就会出现调用断掉
因此我们得加for，手动遍历循环这些middleware，避免middleware中断掉的情况

这里存在一个容易疑惑的点是：
假设T1 开始执行`Next()`，`index`变为0，然后进入for，for内看着似乎会遍历所有的middleware，假设执行到T2 某个middleware内也包含`Next()`调用，此时再次进入`Next()`，又开启了一次for循环，那会不会导致部分middleware重复执行呢？

我们做个假设，现在有4个middleware A、B、C、D,那加上原本的路由处理函数h，我们就得到了5个handler。
从最极端的假设开始，A、B、C、D内部都调用了`Next()`，它们的结构如下：
```go
func(ctx *gee.Context) {
    before()
    ctx.Next()
	after()
}
```
首先路由请求，进入`Next()`，此时`index++`变为`0`，然后开始for循环第一次迭代，for循环先执行`middleware[0]`也即A，A内调用`Next()`，`index++`更新为`1`，然后开启第二个for循环第一次迭代，由于此时`index == 1`，因此执行B，重复如此，直到D，D内调用`Next()`的时候，`index`被更新为4，因此本次for执行的是`middleware[4]`也即路由处理函数h，
```go
//index为4 s为5
for ; c.index < s; c.index++ {
	//c.middleware[c.index]是h
	c.middleware[c.index](c)
}
```
h执行完成后，for循环迭代`c.index++`，此时`index`被更新为`5`，不满足for循环退出，然后回到它的调用点，也即回到D的`Next()`执行后，D执行`after()`，执行完后，for循环迭代`c.index++`，此时`index`被更新为`6`，同样不满足迭代返回到调用点，此时到C的`after()`，重复迭代一直到A的`after()`，最后整个middleware执行完毕。

所以上面这种情况，虽然会多次进入for循环，但每个for其实只会迭代一次，不会重复执行middleware

再考虑另一个情况，假设B和C内都无`Next()`调用，我们再分析会发生什么：

首先路由请求，进入`Next()`，此时`index++`变为`0`，然后开始for循环第一次迭代，for循环先执行`middleware[0]`也即A，A内调用`Next()`，`index++`更新为`1`，然后开启第二个for循环第一次迭代，由于此时`index == 1`，因此执行B，B内无`Next()`因此B执行完成后控制权回到for，for循环迭代`c.index++`，此时`index`被更新为`2`，然后**满足**for继续迭代，此时执行C，同样C内无`Next()`，执行完后，for循环迭代`c.index++`，`index`被更新为`3`，然后**满足**for继续迭代，开始执行D，
D内包含`Next()`调用`Next()`的时候`index`被更新为`4`，执行f，f执行完成，`index`被更新为`5`，不满足继续迭代，控制权回到D的`after()`,D的`after()`执行完成，`index`被更新为`6`，
不满足继续迭代，控制权此时回到**A**，A执行`after()`,`index`被更新为`6`，不满足继续迭代，整个流程执行完成。

可以看到，我们总结下就是，对于中间件中不包含`Next()`调用的，是由for循环的迭代`c.index++`来实现调用下一个middleware的，而中间件如果包含`Next()`，则是通过进入Next一开始的那行`c.index++`实现调用下一middleware的


