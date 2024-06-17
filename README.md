# Matcha

Matcha is the first truly comprehensive, web-based personal finance tracker and budgeting tool.

[![Go](https://github.com/matcha-devs/matcha/actions/workflows/go.yml/badge.svg)](
https://github.com/matcha-devs/matcha/actions/workflows/go.yml)

## Goal features:

* Reliable connections and performance
* Net worth trajectory
* Optimal paycheck allocator
* AI transaction categorizing
* Rewards currency value integration
* Investment balancing
* True value of different tax bucket holdings
* Subscriptions summaries
* Global accounts integration

## Usage

Access our service and create your account at https://www.[domain].com

## Contributing

Pull requests are welcome. For major changes, please open an issue first
to discuss what you would like to change.
Please make sure to update tests as appropriate.

1. Clone this repo
2. To prevent tracking/modifying the run config file, **run**:
   ```git update-index --skip-worktree .idea/runConfigurations/go_build_github_com_matcha_devs_matcha.xml```
3. To enable your databaase, **add** your SQL password to the ```MYSQL_PASSWORD```
   **environment variable** in the **config file**.
4. To import dependencies, **run**: ```go mod tidy```.
5. To run tests, **run**: ```go test ./...```.

## Dependencies

* Web service developed in **Go**
* Database implemented with **MySQL**

## Authors

* [Seoyoung Cho](https://github.com/seoyoungcho213) **co-owner**
* [Carlos Cotera](https://github.com/carlosacj55) **co-owner**
* [Ali A Shah](https://github.com/alishah634)
* [Andrea Goh](https://github.com/andreag0101)
* [Faaiz Memon](https://github.com/faaizmemonpurdue)

## License

[LICENSE.md](LICENSE.md)
