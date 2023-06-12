# go-project
GO高并发实例、GO Web实例

## GO高并发

- (1) 双层队列

由Task队列、Work队列组成。如果Task产生的速度超过Work处理速度，Task会累积并等待。

通过两层队列，避免在Work执行缓慢时(可能是系统原因引起，并不是Task本身处理速度的原因)，过快地拒绝Task。

- (2) 无锁

双层队列都通过 chan 实现，阻塞和等待由 chan 原生机制实现。

- (3) Work 的实现

Work也设计成 chan Task，简化调度器的设计。调度器就简化为把 Task 从一个 chan 移动到另一个 chan。



## GO Web

- (1) 路由

基于 go-martini 实现路由机制，go-martini 的灵活性在于可以任意定义 web controller 的 Input 对象类型；在使用机制上类似于 Spring。

- (2) ORM

一个简化实现的通用 DAL ，但没有实现完整的 ORM ，在 DAL 基础上再进一步，可以实现 ORM ，虽然不是非常复杂，但需要很多时间。

- (3) 分布式唯一的 Sequence Id

实现了一个简单的 Sequence Id Web 服务，对于需要在多个服务之间保证编号唯一性的需要，通过 SequenceId 来完成。

- (4) 支付(支付宝 & 微信支付)

实现了支持支付宝和微信支付的接口，如果需要应用在其他项目，完全可以参考使用。


