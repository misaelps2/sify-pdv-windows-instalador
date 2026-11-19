# Sify PDV para Windows

O instalador cria atalhos no menu Iniciar e na área de trabalho. O atalho abre **https://pdv.sifysistemas.com/** em uma janela de aplicativo do Microsoft Edge. O atalho exibe o ícone da Sify Comércio e mostra uma tela de abertura enquanto a janela inicia. A janela abre em tela cheia, sem a barra superior. Use **Alt+F4** para fechar e **Alt+Tab** para alternar para outro aplicativo. Dentro do PDV, F5, F11, F12 e os atalhos comuns de recarga e inspeção ficam bloqueados. Isso limita o uso desses atalhos no caixa. O código do site entregue ao navegador continua acessível por outros meios.

A página e as atualizações vêm do mesmo endereço usado no navegador. A instalação é por usuário, sem pedir acesso de administrador.

Requisitos: Windows 10 ou 11 e Microsoft Edge (versão estável). Abra o PDV pela primeira vez com internet e aguarde o carregamento e a ativação do caixa. Depois, o funcionamento offline depende do cache e das funções offline do próprio PDV. O Sify Print continua sendo instalado à parte para emissão de cupons.

Feche a janela antiga do Sify PDV antes de atualizar o instalador. A atualização reutiliza o perfil existente.

## Dados e backup

O atalho abre o PDV com o perfil exclusivo em `%LOCALAPPDATA%\Sify PDV\Perfil`. O Chrome e o perfil de uso comum do Edge armazenam os dados em outros diretórios. Reinstalações e desinstalações do Sify PDV preservam a pasta do perfil.

O perfil exclusivo é armazenamento local. Para recuperar dados após perda do disco ou troca de computador, feche o caixa e use **Configurações → Backup e Restauração → Exportar Backup**. Guarde o arquivo criptografado `.pdvbak` e a senha em locais separados. Conserve uma cópia do arquivo fora do computador do caixa, por exemplo em um pendrive. A exportação manual pode ser usada sem contratar ou ativar sincronização em nuvem.

### Migrar um caixa em uso

1. Feche o caixa na janela antiga e exporte o backup por **Configurações → Backup e Restauração**.
2. Instale e abra o Sify PDV com internet. Prepare a ativação se o PDV identificá-lo como um novo perfil.
3. Abra **Configurações → Backup e Restauração → Importar Backup** no aplicativo novo. Informe a senha usada para proteger o arquivo.
4. Confira o histórico de vendas, os cadastros e o caixa antes de iniciar uma venda. Preserve a janela antiga e o arquivo de backup até concluir a conferência.

O perfil novo começa separado. Configurações vinculadas à sessão ou ao dispositivo podem exigir novo login ou configuração. Recursos exclusivos da ponte Electron, como a integração Mercado Pago Pix dessa variante do PDV, dependem do aplicativo Electron que fornece a ponte.

## Compilar

Requisitos: Go 1.25+ e NSIS 3.x. O go-winres gera o recurso de ícone na compilação. O lançador é compilado para Windows x86 e também roda em Windows x64. O GitHub Actions empacota o instalador e verifica a instalação, a abertura no Edge com perfil dedicado e a preservação da pasta do perfil na desinstalação.

No Windows, com Go e NSIS disponíveis no PATH:

```powershell
New-Item -ItemType Directory -Force dist
go install github.com/tc-hib/go-winres@v0.3.3
$winres = Join-Path (go env GOPATH) 'bin\go-winres.exe'
& $winres simply --arch 386 --icon assets/sify-pdv.ico --manifest none --product-name "Sify PDV" --product-version 1.2.0 --file-version 1.2.0 --file-description "Sify Comércio para Windows" --original-filename SifyPDV.exe
$env:GOOS = "windows"
$env:GOARCH = "386"
go build -trimpath -ldflags="-s -w -H windowsgui" -o dist/SifyPDV.exe .
makensis installer.nsi
```

O arquivo pronto é `dist\SifyPDV-Setup.exe`. O instalador aceita `/S` para instalação silenciosa. A versão inicial é de homologação e está sem assinatura Authenticode.
