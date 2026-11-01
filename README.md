# gitnomad

Descarga todos tus repos de GitHub (públicos y privados) en cualquier equipo
con Windows, macOS o Linux, para trabajarlos en local.

> Estado: **fase 1** (motor + CLI: listar y clonar lo que falta).
> Próximas fases: sincronización completa, interfaz gráfica (Wails) y releases.

Web con guía de instalación: **https://rodrigoquezadab.github.io/gitnomad**

## Instalar

```powershell
# Windows (PowerShell)
irm https://rodrigoquezadab.github.io/gitnomad/install.ps1 | iex
```

```sh
# macOS / Linux
curl -fsSL https://rodrigoquezadab.github.io/gitnomad/install.sh | sh
```

Los instaladores detectan sistema y arquitectura, descargan el ejecutable de la
última release, verifican su SHA-256 y lo agregan al PATH. Aceptan `--dir` / `-Dir`
para guardar de una vez la carpeta de repos.

## Requisitos (una vez por equipo)

| | Windows | macOS | Linux |
|---|---|---|---|
| git | `winget install --id Git.Git -e` | `xcode-select --install` | `sudo apt install git` |
| GitHub CLI | `winget install --id GitHub.cli -e` | `brew install gh` | [instrucciones](https://github.com/cli/cli#installation) |

```sh
gh auth login          # GitHub.com → HTTPS → navegador
```

## Uso

```sh
gitnomad doctor --fix                      # revisa todo y arregla lo que pueda
gitnomad config init --dir ~/dev/github    # carpeta de repos en este equipo
gitnomad list                              # qué repos hay y cuáles faltan
gitnomad clone --dry-run                   # ver qué se clonaría
gitnomad clone                             # clonar lo que falta
```

En Windows la carpeta puede ser, por ejemplo, `--dir C:\Dev\github`.

### Comandos

| comando | qué hace |
|---|---|
| `gitnomad doctor` | revisa que git, gh y la sesión de GitHub estén listos |
| `gitnomad doctor --fix` | además arregla lo que pueda (por ejemplo, configurar git para usar las credenciales de gh) |
| `gitnomad config show` | muestra la configuración y dónde está el archivo |
| `gitnomad config init [--dir RUTA] [--pick] [--force]` | crea la configuración de este equipo; `--force` la reemplaza si ya existe |
| `gitnomad config set <clave> <valor>` | cambia una opción (ver [Configuración](#configuración)) |
| `gitnomad config set base_dir --pick` | elige la carpeta de repos en una ventana |
| `gitnomad list [--json]` | lista tus repos: descargado, falta, conflicto u omitido |
| `gitnomad clone` | clona los repos que faltan |
| `gitnomad clone --dry-run` | muestra qué clonaría, sin hacer nada |
| `gitnomad clone --parallel N` | clona N repos a la vez (solo esta vez) |
| `gitnomad clone --dir RUTA` | clona en otra carpeta (solo esta vez) |
| `gitnomad version` | muestra la versión |

### Elegir la carpeta en una ventana

En vez de escribir la ruta, `--pick` abre el selector de carpetas de tu sistema
y guarda la que elijas:

```sh
gitnomad config set base_dir --pick    # cambiar la carpeta
gitnomad config init --pick            # o al crear la configuración
```

- **Windows:** usa PowerShell (viene con el sistema).
- **macOS:** usa `osascript` (viene con el sistema).
- **Linux:** necesita `zenity` o `kdialog` y un escritorio gráfico. En un
  servidor sin escritorio, escribe la ruta: `gitnomad config set base_dir RUTA`.

La ventana abre en la carpeta actual (o en la más cercana que exista). Si la
cierras sin elegir, no se cambia nada.

En la web, el asistente pide en cambio **pegar la ruta**: así los comandos
salen listos para ejecutar y ves exactamente dónde se descargarán los repos.
Si pegas la ruta entre comillas (como la copia Windows con «Copiar como ruta de
acceso»), el asistente las quita.

> Puedes apuntar a una carpeta que ya tenga otros proyectos (por ejemplo
> `C:\Users\tú\Desktop\Code`). Cada repo se clona en su propia subcarpeta, los
> que ya estén ahí se reconocen como descargados y nada existente se modifica.

### Configuración

`gitnomad config set <clave> <valor>`:

| clave | por defecto | qué hace |
|---|---|---|
| `base_dir` | `~/dev/github` | carpeta donde se clona cada repo |
| `parallel` | `4` | repos a la vez |
| `include_forks` | `true` | incluir forks |
| `include_archived` | `false` | incluir repos archivados |
| `exclude` | (vacío) | patrones separados por coma, ej. `test-*,viejo` |

El archivo se guarda en la carpeta de configuración del sistema
(`%APPDATA%\gitnomad`, `~/Library/Application Support/gitnomad`, `~/.config/gitnomad`).

## Seguridad

- Nunca borra ni sobrescribe nada local. Si existe una carpeta con el nombre
  de un repo pero no es ese repo, se marca como **conflicto** y no se toca.
- No guarda tokens: usa la sesión de `gh`.
- No pide claves por consola (no se queda colgado).

## Compilar

Requiere Go 1.22+. En Windows: `winget install --id GoLang.Go -e`.

```sh
go test ./...
go build -o dist/ ./cmd/gitnomad
# cruzado:
GOOS=windows GOARCH=amd64 go build -o dist/gitnomad.exe ./cmd/gitnomad
GOOS=darwin  GOARCH=arm64 go build -o dist/gitnomad-macos ./cmd/gitnomad
```

## Publicar

1. **Sitio web** (una sola vez): en GitHub, *Settings → Pages → Build and deployment →
   Source: GitHub Actions*. Cada push a `main` que toque `web/` o `scripts/install.*`
   publica el sitio con los instaladores.
2. **Versión nueva**: `git tag v0.1.0 && git push origin v0.1.0`. El workflow
   `release.yml` corre los tests, compila los 6 ejecutables con `scripts/build.sh`
   y crea la release con `checksums.txt`. Los instaladores siempre bajan la última.

Si haces un fork, cambia `rodrigoquezadab/gitnomad` en `web/app.js` (`CONFIG`) y en
los dos instaladores (o usa la variable `GITNOMAD_REPO`).

## Estructura

```
cmd/gitnomad/       CLI
internal/config/   configuración por equipo
internal/github/   lista de repos vía gh
internal/gitops/   ejecutar git, clonar, comparar URLs
internal/syncer/   motor: filtrar, inspeccionar, clonar en paralelo
internal/doctor/   chequeos del equipo
internal/picker/   selector de carpetas nativo (--pick)
scripts/           instaladores (install.sh, install.ps1) y build.sh
web/               sitio (GitHub Pages): guía bilingüe con asistente
.github/workflows/ ci.yml (tests en 3 SO), release.yml, pages.yml
```

## Desarrollo local

**Probar el CLI sin instalarlo.** Desde la carpeta del repo:

```sh
go run ./cmd/gitnomad clone --dry-run
```

El instalador de una línea (`irm ... | iex` o `curl ... | sh`) siempre baja la
**última release publicada**, no tu código local. Para probar cambios que aún
no están en una release, usa `go run` o compila con `go build`.

**Sin tocar tu configuración real.** Apunta a un archivo temporal con
`GITNOMAD_CONFIG`:

```powershell
$env:GITNOMAD_CONFIG = "$env:TEMP\gitnomad-prueba.json"   # PowerShell
```

```sh
export GITNOMAD_CONFIG=/tmp/gitnomad-prueba.json          # macOS / Linux
```

**`go` no se reconoce tras instalarlo (Windows).** Las terminales que ya estaban
abiertas no ven el PATH nuevo. Abre una nueva o, en la actual:

```powershell
$env:Path = "C:\Program Files\Go\bin;$env:Path"
```

**Ver el sitio en local** (se recarga al guardar):

```sh
npx live-server web --port=5500
```

O, sin recarga automática, `python3 -m http.server -d web 8000`.
