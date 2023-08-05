Ferramenta simples para procurar IPs, sejam publicos ou privados, de todas as instâncias EC2 do perfil passado através da CLI.

- **--ip_type="private" ou --ip_type="public"**

```
go run main.go --profile="default" --region="us-east-1" --ip_type="private"
```
Irá não só printar no terminal os IPs públicos ou privados, como também irá gerar um arquivo de logs. Ideal para criar inventário para o Ansible.
