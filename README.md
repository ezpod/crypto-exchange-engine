# crypto-exchange-engine

加密资产交易撮合引擎的POC实现代码。

交易撮合引擎（Matching/Trading Engine），顾名思义是用来撮合交易的软件，广泛地应用在
金融、证券、加密货币交易等领域。交易引擎负责管理加密资产市场中所有的开口订单（Open Orders），
并在发现匹配的订单对（Trading Pair）时自动执行交易。本文将首先介绍有关加密资产交易撮合
引擎的基本概念，例如委托单、交易委托账本等，然后使用Golang实现一个原理性的撮合引擎。如果你正在
考虑实现类似交易所（Exchange）这样的产品，相信本文会对你有很大的帮助。

完整说明：[交易撮合引擎原理与实现](http://blog.hubwiz.com/2019/07/22/build-a-crypto-exchange/)

![](http://blog.hubwiz.com/images/cta-1.png)

如果你想学习区块链并在Blockchain Technologies建立职业生涯，那么请查看我们分享的一些以太坊、比特币、EOS、Fabric等区块链相关的交互式在线编程实战教程：

> - [java以太坊开发教程](http://xc.hubwiz.com/course/5b2b6e82c02e6b6a59171de2?affid=github7878)，主要是针对java和android程序员进行区块链以太坊开发的web3j详解。
> - [python以太坊](http://xc.hubwiz.com/course/5b40462cc02e6b6a59171de4?affid=github7878)，主要是针对python工程师使用web3.py进行区块链以太坊开发的详解。
> - [php以太坊](http://xc.hubwiz.com/course/5b36629bc02e6b6a59171de3?affid=github7878)，主要是介绍使用php进行智能合约开发交互，进行账号创建、交易、转账、代币开发以及过滤器和交易等内容。
> - [以太坊入门教程](http://xc.hubwiz.com/course/5a952991adb3847553d205d1?affid=github7878)，主要介绍智能合约与dapp应用开发，适合入门。
> - [以太坊开发进阶教程](http://xc.hubwiz.com/course/5abbb7acc02e6b6a59171dd6?affid=github7878)，主要是介绍使用node.js、mongodb、区块链、ipfs实现去中心化电商DApp实战，适合进阶。
> - [ERC721以太坊通证实战](http://xc.hubwiz.com/course/5c6ed395070c379b559a813a?affid=github7878)，课程以一个数字艺术品创作与分享DApp的实战开发为主线，深入讲解以太坊非同质化通证的概念、标准与开发方案。内容包含ERC-721标准的自主实现，讲解OpenZeppelin合约代码库二次开发，实战项目采用Truffle，IPFS，实现了通证以及去中心化的通证交易所。
> - [C#以太坊](http://xc.hubwiz.com/course/5b6048c3c02e6b6a59171dee?affid=github7878)，主要讲解如何使用C#开发基于.Net的以太坊应用，包括账户管理、状态与交易、智能合约开发与交互、过滤器和交易等。
> - [java比特币开发教程](http://xc.hubwiz.com/course/5bb35c90c02e6b6a59171df0?affid=github7878)，本课程面向初学者，内容即涵盖比特币的核心概念，例如区块链存储、去中心化共识机制、密钥与脚本、交易与UTXO等，同时也详细讲解如何在Java代码中集成比特币支持功能，例如创建地址、管理钱包、构造裸交易等，是Java工程师不可多得的比特币开发学习课程。
> - [php比特币开发教程](http://xc.hubwiz.com/course/5b9e779ac02e6b6a59171def?affid=github7878)，本课程面向初学者，内容即涵盖比特币的核心概念，例如区块链存储、去中心化共识机制、密钥与脚本、交易与UTXO等，同时也详细讲解如何在Php代码中集成比特币支持功能，例如创建地址、管理钱包、构造裸交易等，是Php工程师不可多得的比特币开发学习课程。
> - [c#比特币开发教程](http://xc.hubwiz.com/course/5c766a59f54a5e207931b5a5?affid=github7878)，本课程面向初学者，内容即涵盖比特币的核心概念，例如区块链存储、去中心化共识机制、密钥与脚本、交易与UTXO等，同时也详细讲解如何在C#代码中集成比特币支持功能，例如创建地址、管理钱包、构造裸交易等，是C#工程师不可多得的比特币开发学习课程。
> - [EOS入门教程](http://xc.hubwiz.com/course/5b52c0a2c02e6b6a59171ded?affid=github7878)，本课程帮助你快速入门EOS区块链去中心化应用的开发，内容涵盖EOS工具链、账户与钱包、发行代币、智能合约开发与部署、使用代码与智能合约交互等核心知识点，最后综合运用各知识点完成一个便签DApp的开发。
> - [深入浅出玩转EOS钱包开发](http://xc.hubwiz.com/course/5c79edcaf697372707791512?affid=github7878)，本课程以手机EOS钱包的完整开发过程为主线，深入学习EOS区块链应用开发，课程内容即涵盖账户、计算资源、智能合约、动作与交易等EOS区块链的核心概念，同时也讲解如何使用eosjs和eosjs-ecc开发包访问EOS区块链，以及如何在React前端应用中集成对EOS区块链的支持。课程内容深入浅出，非常适合前端工程师深入学习EOS区块链应用开发。
> - [Hyperledger Fabric 区块链开发详解](http://xc.hubwiz.com/course/5c9b89f54898e59b7b63430a?affid=github7878)，本课程面向初学者，内容即包含Hyperledger Fabric的身份证书与MSP服务、权限策略、信道配置与启动、链码通信接口等核心概念，也包含Fabric网络设计、nodejs链码与应用开发的操作实践，是Nodejs工程师学习Fabric区块链开发的最佳选择。
> - [Hyperledger Fabric java 区块链开发详解](http://xc.hubwiz.com/course/5c9b89f54898e59b7b63430a?affid=github7878)，课程面向初学者，内容即包含Hyperledger Fabric的身份证书与MSP服务、权限策略、信道配置与启动、链码通信接口等核心概念，也包含Fabric网络设计、java链码与应用开发的操作实践，是java工程师学习Fabric区块链开发的最佳选择。
> - [tendermint区块链开发详解](http://xc.hubwiz.com/course/5bdec63ac02e6b6a59171df3?affid=github7878)，本课程适合希望使用tendermint进行区块链开发的工程师，课程内容即包括tendermint应用开发模型中的核心概念，例如ABCI接口、默克尔树、多版本状态库等，也包括代币发行等丰富的实操代码，是go语言工程师快速入门区块链开发的最佳选择。
