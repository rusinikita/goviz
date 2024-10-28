# goviz
Go repository architecture visualization

![](/docs/top_level.png)
![](/docs/tooltip.png)

# Usage

1. Clone repository
2. Download dependencies
```shell
go mod tidy
```
3. Replace value of `projectDir` and `pathPrefixToRemove` with your local repository values
4. Run app
```shell
go run main.go
```
5. Open `result.html` file in browser