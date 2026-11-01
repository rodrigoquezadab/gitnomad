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

Para elegir la carpeta en el selector de carpetas del sistema en vez de
escribirla, usa `--pick`: `gitnomad config init --pick` o
`gitnomad config set base_dir --pick`. En Linux necesita `zenity` o `kdialog`.

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

Requiere Go 1.22+.

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
scripts/           instaladores (install.sh, install.ps1) y build.sh
web/               sitio (GitHub Pages): guía bilingüe con asistente
.github/workflows/ ci.yml (tests en 3 SO), release.yml, pages.yml
```

Para ver el sitio en local: `python3 -m http.server -d web 8000`.
