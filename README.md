# gitnomad

Descarga **todos tus repos de GitHub** (públicos y privados) en cualquier equipo
con Windows, macOS o Linux, con un solo comando. Úsalo cuando estrenas
computador, cambias de equipo o quieres tener todo tu código en local.

- Clona solo lo que falta: puedes repetirlo cuando quieras.
- Nunca borra ni sobrescribe nada que ya tengas.
- Usa tu sesión de GitHub CLI: no tienes que crear ni copiar tokens.

> Estado: **fase 1** (listar y clonar lo que falta). Próximamente: traer cambios
> de repos ya descargados, interfaz gráfica y repos de organizaciones.

Guía interactiva (arma los comandos por ti): **https://rodrigoquezadab.github.io/gitnomad**

## Contenido

- [Inicio rápido](#inicio-rápido) (10 minutos, una vez por equipo)
- [Comandos](#comandos)
- [Configuración](#configuración)
- [Problemas frecuentes](#problemas-frecuentes)
- [Seguridad](#seguridad)
- [Para desarrolladores](#para-desarrolladores)

---

## Inicio rápido

Sigue los pasos en orden. Después de cada uno se indica **qué deberías ver**;
si no lo ves, revisa [Problemas frecuentes](#problemas-frecuentes) antes de seguir.

> **Regla de oro:** cada vez que instales algo (git, gh, gitnomad o Go),
> **cierra la terminal y abre una nueva**. Las terminales que ya estaban abiertas
> no ven los programas recién instalados y dirán que el comando «no se reconoce».

### 1. Instala git y GitHub CLI

**Windows** (en PowerShell):

```powershell
winget install --id Git.Git -e
winget install --id GitHub.cli -e
```

**macOS** (requiere [Homebrew](https://brew.sh)):

```sh
brew install git gh
```

**Linux** (Debian/Ubuntu; para otras distribuciones mira
[las instrucciones de gh](https://github.com/cli/cli#installation)):

```sh
sudo apt update && sudo apt install -y git gh
```

Cierra la terminal y abre una nueva. **Qué deberías ver:**

```sh
git --version     # git version 2.x
gh --version      # gh version 2.x
```

### 2. Inicia sesión en GitHub

```sh
gh auth login
```

Responde así a las preguntas:

| pregunta | respuesta |
|---|---|
| Where do you use GitHub? | `GitHub.com` |
| Preferred protocol for Git operations? | `HTTPS` |
| Authenticate Git with your GitHub credentials? | `Yes` |
| How would you like to authenticate? | `Login with a web browser` |

Copia el código que aparece, pulsa Enter y pégalo en el navegador.
**Qué deberías ver:**

```sh
gh auth status    # ✓ Logged in to github.com account TU-USUARIO
```

### 3. Instala gitnomad

**Windows** (PowerShell):

```powershell
irm https://rodrigoquezadab.github.io/gitnomad/install.ps1 | iex
```

**macOS / Linux:**

```sh
curl -fsSL https://rodrigoquezadab.github.io/gitnomad/install.sh | sh
```

El instalador descarga el ejecutable de la **última versión publicada**, verifica
su SHA-256 y lo agrega a tu PATH (Windows: `%LOCALAPPDATA%\Programs\gitnomad`;
macOS/Linux: `~/.local/bin`).

Cierra la terminal y abre una nueva. **Qué deberías ver:**

```sh
gitnomad version  # gitnomad v0.x.x (windows/amd64)
```

### 4. Elige la carpeta donde se descargarán los repos

Copia la ruta completa de la carpeta y guárdala:

```powershell
gitnomad config set base_dir "C:\Users\TU-USUARIO\Code"      # Windows
```

```sh
gitnomad config set base_dir ~/code                          # macOS / Linux
```

Cómo copiar la ruta:

- **Windows:** abre la carpeta en el Explorador de archivos, haz clic en la barra
  de direcciones y copia (Ctrl+C).
- **macOS:** selecciona la carpeta en Finder y pulsa Opción+Cmd+C.
- **Linux:** en el explorador de archivos pulsa Ctrl+L y copia, o ejecuta `pwd`
  dentro de la carpeta.

Usa comillas si la ruta tiene espacios. Si la carpeta no existe, se crea al clonar.
Puede ser una carpeta que ya tenga otros proyectos: cada repo va a su propia
subcarpeta y lo que ya existe no se toca.

**Qué deberías ver:** `Guardado.` y una tabla con `base_dir` apuntando a tu
carpeta. Este paso **solo guarda la configuración**: todavía no descarga nada.

> ¿Prefieres elegirla en una ventana? `gitnomad config set base_dir --pick` abre
> el selector de carpetas del sistema (ver [Elegir la carpeta en una
> ventana](#elegir-la-carpeta-en-una-ventana)).

### 5. Revisa y clona

```sh
gitnomad doctor --fix      # revisa que todo esté listo y arregla lo que pueda
gitnomad clone --dry-run   # muestra qué se clonaría y dónde, sin descargar nada
gitnomad clone             # descarga los repos que faltan
```

**Qué deberías ver:** `doctor` termina con `Todo listo.`; `--dry-run` lista cada
repo con la carpeta de destino; `clone` muestra `[1/N] ok ...` por cada repo y
termina con `Listo en Ns: N clonados, 0 con error.`

¡Ya está! En otro equipo repite estos pasos. Cuando crees repos nuevos en GitHub,
vuelve a ejecutar `gitnomad clone`: solo descargará los que falten.

---

## Comandos

| comando | qué hace |
|---|---|
| `gitnomad doctor` | revisa que git, gh y la sesión de GitHub estén listos |
| `gitnomad doctor --fix` | además arregla lo que pueda (por ejemplo, configura git para usar las credenciales de gh) |
| `gitnomad config show` | muestra la configuración y dónde está el archivo |
| `gitnomad config set <clave> <valor>` | cambia una opción (ver [Configuración](#configuración)) |
| `gitnomad config set base_dir --pick` | elige la carpeta de repos en una ventana |
| `gitnomad config init [--dir RUTA] [--pick] [--force]` | crea la configuración desde cero; `--force` reemplaza una existente |
| `gitnomad list` | lista tus repos: `descargado`, `falta`, `conflicto` u `omitido` (`--json` para scripts) |
| `gitnomad clone --dry-run` | muestra qué clonaría, sin hacer nada |
| `gitnomad clone` | clona los repos que faltan |
| `gitnomad clone --parallel N` | clona N repos a la vez (solo esta vez) |
| `gitnomad clone --dir RUTA` | clona en otra carpeta (solo esta vez) |
| `gitnomad version` | muestra la versión instalada |

## Configuración

Cada equipo tiene su propia configuración, así que en cada uno puedes usar una
carpeta distinta. Cámbiala con `gitnomad config set <clave> <valor>`:

| clave | por defecto | qué hace |
|---|---|---|
| `base_dir` | `~/dev/github` | carpeta donde se clona cada repo |
| `parallel` | `4` | repos que se clonan a la vez |
| `include_forks` | `true` | incluir los repos que son fork |
| `include_archived` | `false` | incluir los repos archivados |
| `exclude` | (vacío) | patrones de nombres a ignorar, separados por coma: `"test-*,viejo"` |

El archivo se guarda en:

- Windows: `%APPDATA%\gitnomad\config.json`
- macOS: `~/Library/Application Support/gitnomad/config.json`
- Linux: `~/.config/gitnomad/config.json`

### Elegir la carpeta en una ventana

`--pick` abre el selector de carpetas de tu sistema y guarda la que elijas:

```sh
gitnomad config set base_dir --pick
```

- **Windows y macOS:** funciona sin instalar nada.
- **Linux:** necesita `zenity` o `kdialog` y un escritorio gráfico. En un
  servidor sin escritorio, escribe la ruta como en el [paso 4](#4-elige-la-carpeta-donde-se-descargarán-los-repos).

La ventana abre en la carpeta actual. Si la cierras sin elegir, no cambia nada.

> `--pick` está disponible a partir de **v0.2.0**. Comprueba tu versión con
> `gitnomad version`; si es anterior, reinstala con el comando del paso 3.

---

## Problemas frecuentes

**«gitnomad / gh / git / go no se reconoce como nombre de un cmdlet...»**
La terminal se abrió antes de instalar el programa. Ciérrala y abre una nueva
(si usas una terminal integrada en un editor o app, reinicia esa app). Si
sigue fallando en Windows, ejecuta de nuevo el instalador del paso 3.

**Ejecuté `gitnomad config set base_dir ...` y «no pasó nada».**
Es lo esperado: ese comando solo guarda la carpeta en la configuración. Los repos
se descargan con `gitnomad clone` (prueba antes `gitnomad clone --dry-run`).

**`--pick` da `error: base_dir debe ser una ruta absoluta (recibido: "--pick")`.**
Tienes una versión anterior a v0.2.0, que no conoce `--pick` (no se guardó nada).
Reinstala (paso 3), comprueba con `gitnomad version` y vuelve a ejecutarlo, o
escribe la ruta como en el [paso 4](#4-elige-la-carpeta-donde-se-descargarán-los-repos).

**¿Por qué el instalador no trae los últimos cambios del repositorio?**
El instalador descarga la última **versión publicada** (release), no el código
de la rama `main`. Para usar cambios aún no publicados, mira
[Para desarrolladores](#para-desarrolladores).

**Un repo aparece como `conflicto`.**
En tu carpeta ya existe una subcarpeta con el nombre de ese repo, pero no es una
copia de él. gitnomad no la toca: renómbrala o muévela y vuelve a ejecutar
`gitnomad clone`.

**Quiero clonar en una carpeta que ya tiene otros proyectos.**
Puedes. Los repos que ya estén ahí se reconocen como descargados, el resto se
clona en subcarpetas nuevas y nada existente se modifica. Revisa antes con
`gitnomad clone --dry-run`.

**Errores de autenticación o no aparecen mis repos privados.**
Ejecuta `gh auth status`. Si no dice `Logged in`, repite `gh auth login`
(paso 2). Luego ejecuta `gitnomad doctor --fix`.

**No aparecen los repos de mi organización.**
Por ahora gitnomad descarga solo los repos de tu cuenta personal.

---

## Seguridad

- Nunca borra ni sobrescribe nada local. Si existe una carpeta con el nombre de
  un repo pero no es ese repo, se marca como **conflicto** y no se toca.
- No guarda tokens: usa la sesión oficial de `gh`.
- No pide claves por consola, así que nunca se queda colgado esperando.
- Los instaladores verifican la suma SHA-256 de cada descarga.

---

## Para desarrolladores

Esta sección es para quien quiere modificar gitnomad o probar cambios que aún no
están en una versión publicada. **Si solo quieres usar la herramienta, no la
necesitas**: sigue el [Inicio rápido](#inicio-rápido).

### Preparar el entorno

1. Instala **Go 1.22 o superior**, además de git y gh (pasos 1 y 2):

   ```powershell
   winget install --id GoLang.Go -e      # Windows
   ```

   ```sh
   brew install go                       # macOS
   ```

   Cierra la terminal y abre una nueva; `go version` debe responder.

2. Clona el repositorio y entra en él:

   ```sh
   git clone https://github.com/rodrigoquezadab/gitnomad.git
   cd gitnomad
   ```

3. **Solo en Windows:** evita que git marque los scripts `.sh` como modificados
   por los permisos de archivo, que Windows no maneja:

   ```sh
   git config core.fileMode false
   ```

4. Comprueba que todo compila y los tests pasan:

   ```sh
   go test ./...
   ```

### Probar tus cambios

Ejecuta el CLI directamente desde el código, sin instalarlo:

```sh
go run ./cmd/gitnomad version
go run ./cmd/gitnomad clone --dry-run
```

Para no tocar tu configuración real mientras pruebas, apunta a un archivo
temporal (dura mientras la terminal esté abierta):

```powershell
$env:GITNOMAD_CONFIG = "$env:TEMP\gitnomad-prueba.json"   # PowerShell
```

```sh
export GITNOMAD_CONFIG=/tmp/gitnomad-prueba.json          # macOS / Linux
```

Para generar el ejecutable:

```sh
go build -o dist/ ./cmd/gitnomad
```

### Ver el sitio web en local

El sitio está en `web/` (HTML, CSS y JS sin dependencias). Con recarga automática
al guardar:

```sh
npx live-server web --port=5500
```

O sin Node: `python3 -m http.server -d web 8000`.

### Estructura

```
cmd/gitnomad/       CLI (punto de entrada y comandos)
internal/config/    configuración por equipo
internal/github/    lista de repos vía gh
internal/gitops/    ejecutar git, clonar, comparar URLs
internal/syncer/    motor: filtrar, inspeccionar, clonar en paralelo
internal/doctor/    chequeos del equipo
internal/picker/    selector de carpetas nativo (--pick)
scripts/            instaladores (install.sh, install.ps1) y build.sh
web/                sitio (GitHub Pages): guía bilingüe con asistente
.github/workflows/  ci.yml (tests en 3 SO), release.yml, pages.yml
```

### Publicar

1. **Sitio web** (una sola vez): en GitHub, *Settings → Pages → Build and
   deployment → Source: GitHub Actions*. Cada push a `main` que toque `web/` o
   `scripts/install.*` publica el sitio con los instaladores.
2. **Versión nueva:** con los cambios ya en `main`, crea una etiqueta:

   ```sh
   git tag v0.2.0
   git push origin v0.2.0
   ```

   El workflow `release.yml` corre los tests, compila los 6 ejecutables con
   `scripts/build.sh` y crea la release con `checksums.txt`. Desde ese momento
   los instaladores descargan esa versión.

Si haces un fork, cambia `rodrigoquezadab/gitnomad` en `web/app.js` (`CONFIG`) y
en los dos instaladores, o usa la variable de entorno `GITNOMAD_REPO`.
