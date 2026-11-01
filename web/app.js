/* gitnomad — sitio web: idioma, asistente de instalación y botones de copiar. */
(function () {
  'use strict';

  // ---- Configuración del proyecto (cambia esto si haces un fork) ----
  var CONFIG = {
    repo: 'rodrigoquezadab/gitnomad',
    site: 'https://rodrigoquezadab.github.io/gitnomad'
  };
  var REPO_URL = 'https://github.com/' + CONFIG.repo;
  var RELEASES_URL = REPO_URL + '/releases/latest';
  var GH_INSTALL_URL = 'https://github.com/cli/cli#installation';

  var DEFAULT_DIR = { windows: 'C:\\Dev\\github', macos: '~/dev/github', linux: '~/dev/github' };

  // ---- Textos ----
  var T = {
    es: {
      'meta.description': 'Descarga todos tus repos de GitHub, públicos y privados, en Windows, macOS o Linux con un solo comando.',
      'nav.skip': 'Saltar a la instalación',
      'nav.how': 'Cómo funciona',
      'nav.start': 'Empezar',
      'nav.cmds': 'Comandos',
      'nav.faq': 'Preguntas',
      'hero.title': 'Todos tus repos de GitHub, en cada equipo.',
      'hero.lead': 'gitnomad descarga tus repositorios públicos y privados en la carpeta que elijas, con un solo comando. Usa tu sesión de GitHub CLI, así que no tienes que copiar tokens, y nunca pisa tu trabajo local.',
      'hero.cta': 'Empezar ahora',
      'hero.repo': 'Ver código en GitHub',
      'how.title': 'Tres pasos, una vez por equipo',
      'how.1.t': 'Vincula tu cuenta',
      'how.1.d': 'Inicias sesión con GitHub CLI desde el navegador. gitnomad usa esa sesión para ver también tus repos privados.',
      'how.2.t': 'Elige la carpeta',
      'how.2.d': 'Cada equipo guarda los repos donde tú quieras: <code>C:\\Dev\\github</code> en Windows, <code>~/dev/github</code> en Mac o Linux.',
      'how.3.t': 'Clona todo',
      'how.3.d': 'gitnomad compara tu cuenta con la carpeta y clona solo lo que falta. Puedes repetirlo cuando quieras.',
      'start.title': 'Instalación guiada',
      'start.lead': 'Elige tu sistema y la carpeta. Los comandos se arman solos: cópialos en orden.',
      'wiz.os': 'Tu sistema',
      'wiz.distro': 'Distribución',
      'wiz.dir': 'Carpeta donde se clonarán los repos',
      'wiz.dir.hint': 'Pega la ruta completa: los comandos de abajo se arman con ella. Se crea si no existe.',
      'wiz.dir.how.windows': 'Para copiarla: abre la carpeta en el Explorador de archivos, haz clic en la barra de direcciones y copia (Ctrl+C). También sirve clic derecho → «Copiar como ruta de acceso».',
      'wiz.dir.how.macos': 'Para copiarla: en Finder selecciona la carpeta y pulsa Opción+Cmd+C (o clic derecho con Opción → «Copiar como nombre de ruta»). Puedes usar <code>~</code> para tu carpeta personal.',
      'wiz.dir.how.linux': 'Para copiarla: en tu explorador de archivos pulsa Ctrl+L y copia la ruta, o escribe <code>pwd</code> en una terminal dentro de la carpeta. Puedes usar <code>~</code> para tu carpeta personal.',
      'wiz.options': 'Opciones',
      'wiz.forks': 'Incluir forks',
      'wiz.archived': 'Incluir repos archivados',
      'wiz.warn.empty': 'Escribe una carpeta para continuar.',
      'wiz.warn.relative': 'Usa una ruta completa, por ejemplo {ex}.',
      's1.t': 'Instala git y GitHub CLI',
      's1.win': 'Abre <strong>PowerShell</strong> y ejecuta (si ya los tienes, sáltate este paso):',
      's1.win.after': 'Cierra PowerShell y ábrelo de nuevo para que reconozca los comandos nuevos.',
      's1.mac': 'Con <a href="https://brew.sh">Homebrew</a>, en la app <strong>Terminal</strong>:',
      's1.linux': 'Con el gestor de paquetes de tu distribución:',
      's1.linux.after': 'Si tu distribución no trae <code>gh</code> o trae una versión antigua, usa la <a href="{url}">guía oficial de GitHub CLI</a>.',
      's2.t': 'Vincula tu cuenta de GitHub',
      's2.d': 'Inicia sesión una sola vez en este equipo:',
      's2.answers': 'GitHub CLI te hará unas preguntas. Responde así:',
      's2.q1': 'Where do you use GitHub?',
      's2.q2': 'What is your preferred protocol for Git operations on this host?',
      's2.q3': 'Authenticate Git with your GitHub credentials?',
      's2.q4': 'How would you like to authenticate GitHub CLI?',
      's2.browser': 'Copia el código de 8 caracteres que aparece, presiona Enter y pégalo en la página de GitHub que se abre. Luego autoriza el acceso.',
      's2.verify': 'Comprueba que quedó vinculada (debe mostrar tu usuario):',
      's2.note': 'GitHub CLI guarda la sesión en tu equipo. gitnomad no lee ni guarda tu token: le pide a <code>gh</code> que se autentique cuando hace falta.',
      's3.t': 'Instala gitnomad',
      's3.win': 'En PowerShell:',
      's3.unix': 'En la terminal:',
      's3.unix.after': 'Si después dice que no encuentra <code>gitnomad</code>, abre una terminal nueva.',
      's3.manual': '¿Prefieres descargarlo a mano? Está en <a href="{url}">Releases</a>.',
      's4.t': 'Define la carpeta de tus repos',
      's4.d': 'Guarda la carpeta en la configuración de este equipo:',
      's4.opts': 'Opciones que elegiste:',
      's4.note': 'Cada equipo tiene su propia configuración, así que en otro dispositivo puedes usar otra ruta.',
      's5.t': 'Revisa y clona',
      's5.doctor': 'Revisa que todo esté listo (y arregla lo que se pueda):',
      's5.list': 'Mira qué repos tienes y cuáles faltan:',
      's5.dry': 'Simula sin descargar nada:',
      's5.clone': 'Clona lo que falta:',
      's5.after': 'En otro equipo repite estos pasos. Cuando crees repos nuevos, vuelve a correr <code>gitnomad clone</code>.',
      'copy': 'Copiar',
      'copied': 'Copiado',
      'copy.toast': 'Comando copiado',
      'copy.fail': 'No se pudo copiar; selecciónalo a mano',
      'safety.title': 'Hecho para no romper nada',
      'safety.1.t': 'No borra ni sobrescribe.',
      'safety.1.d': 'Si ya existe una carpeta con el nombre de un repo y no es ese repo, la marca como conflicto y no la toca.',
      'safety.2.t': 'Sin tokens a la vista.',
      'safety.2.d': 'La autenticación la maneja GitHub CLI con el inicio de sesión oficial de GitHub.',
      'safety.3.t': 'No se queda colgado.',
      'safety.3.d': 'Nunca pide claves por consola. Si un repo falla, lo informa y sigue con el resto.',
      'safety.4.t': 'Instalador verificado.',
      'safety.4.d': 'Cada descarga se comprueba con su suma SHA-256 antes de instalarse.',
      'cmds.title': 'Todos los comandos',
      'cmds.lead': 'Referencia rápida. Las opciones de <code>clone</code> solo valen para esa ejecución; para dejarlas fijas usa <code>gitnomad config set</code>.',
      'cmds.keys': 'Claves de <code>config set</code>: <code>base_dir</code>, <code>parallel</code>, <code>include_forks</code>, <code>include_archived</code> y <code>exclude</code>.',
      cmds: [
        ['gitnomad doctor', 'Revisa que git, gh y tu sesión de GitHub estén listos.'],
        ['gitnomad doctor --fix', 'Además arregla lo que pueda, como configurar git para usar las credenciales de gh.'],
        ['gitnomad config show', 'Muestra tu configuración y dónde está el archivo.'],
        ['gitnomad config init --dir RUTA', 'Crea la configuración de este equipo. Con <code>--pick</code> eliges la carpeta en una ventana; <code>--force</code> reemplaza una existente.'],
        ['gitnomad config set base_dir --pick', 'Abre el selector de carpetas de tu sistema y guarda la que elijas.'],
        ['gitnomad config set <clave> <valor>', 'Cambia una opción, por ejemplo <code>parallel 8</code> o <code>exclude "test-*"</code>.'],
        ['gitnomad list', 'Lista tus repos y si están descargados, faltan, en conflicto u omitidos. Con <code>--json</code> para scripts.'],
        ['gitnomad clone --dry-run', 'Muestra qué se clonaría y dónde, sin descargar nada.'],
        ['gitnomad clone', 'Clona los repos que faltan. Opciones: <code>--parallel N</code> y <code>--dir RUTA</code>.'],
        ['gitnomad version', 'Muestra la versión instalada.']
      ],
      'faq.title': 'Preguntas frecuentes',
      'footer.by': 'hecho por',
      'footer.code': 'Código en GitHub',
      faq: [
        ['¿Qué repos descarga?', 'Los repos públicos y privados de tu cuenta, incluidos los forks. Los archivados quedan fuera salvo que actives la opción. Los repos de organizaciones todavía no se incluyen.'],
        ['¿Actualiza los repos que ya tengo?', 'Por ahora clona los que faltan y no toca los existentes. La sincronización (traer cambios y avisar de commits sin subir) llega en la próxima versión.'],
        ['¿Funciona con varias cuentas de GitHub?', 'Usa la cuenta activa en GitHub CLI. Puedes cambiarla con <code>gh auth switch</code> y apuntar a otra carpeta con <code>gitnomad config set base_dir</code>.'],
        ['¿Dónde se guarda la configuración?', '<p>Windows: <code>%APPDATA%\\gitnomad\\config.json</code></p><p>macOS: <code>~/Library/Application Support/gitnomad/config.json</code></p><p>Linux: <code>~/.config/gitnomad/config.json</code></p>'],
        ['¿Puedo elegir la carpeta en una ventana en vez de escribirla?', 'Sí: <code>gitnomad config set base_dir --pick</code> abre el selector de carpetas de tu sistema. En Windows y macOS funciona sin instalar nada; en Linux necesita <code>zenity</code> o <code>kdialog</code>. Si cierras la ventana sin elegir, no cambia nada. Disponible desde v0.2.0 (mira tu versión con <code>gitnomad version</code>).'],
        ['¿Puedo clonar en una carpeta que ya tiene otros proyectos?', 'Sí. Cada repo va a su propia subcarpeta, los que ya estén ahí se reconocen como descargados y nada existente se modifica. Si una carpeta tiene el nombre de un repo pero es otra cosa, queda como conflicto y no se toca. Usa <code>gitnomad clone --dry-run</code> para ver antes qué pasaría.'],
        ['¿Cómo excluyo algunos repos?', 'Con patrones separados por coma: <code>gitnomad config set exclude "test-*,viejo"</code>.'],
        ['¿Cómo actualizo gitnomad?', 'Vuelve a ejecutar el comando de instalación: reemplaza el ejecutable por la última versión.'],
        ['¿Cómo lo desinstalo?', 'Borra el ejecutable (Windows: <code>%LOCALAPPDATA%\\Programs\\gitnomad</code>; macOS/Linux: <code>~/.local/bin/gitnomad</code>) y la carpeta de configuración. Tus repos clonados no se tocan.']
      ]
    },
    en: {
      'meta.description': 'Download all your GitHub repos, public and private, on Windows, macOS or Linux with a single command.',
      'nav.skip': 'Skip to installation',
      'nav.how': 'How it works',
      'nav.start': 'Get started',
      'nav.cmds': 'Commands',
      'nav.faq': 'FAQ',
      'hero.title': 'All your GitHub repos, on every machine.',
      'hero.lead': 'gitnomad downloads your public and private repositories into the folder you choose, with a single command. It uses your GitHub CLI session, so there are no tokens to copy, and it never overwrites your local work.',
      'hero.cta': 'Get started',
      'hero.repo': 'View code on GitHub',
      'how.title': 'Three steps, once per machine',
      'how.1.t': 'Link your account',
      'how.1.d': 'Sign in with GitHub CLI through your browser. gitnomad uses that session to see your private repos too.',
      'how.2.t': 'Pick a folder',
      'how.2.d': 'Each machine keeps the repos wherever you like: <code>C:\\Dev\\github</code> on Windows, <code>~/dev/github</code> on Mac or Linux.',
      'how.3.t': 'Clone everything',
      'how.3.d': 'gitnomad compares your account with the folder and clones only what is missing. Run it again any time.',
      'start.title': 'Guided setup',
      'start.lead': 'Choose your system and folder. The commands are generated for you: copy them in order.',
      'wiz.os': 'Your system',
      'wiz.distro': 'Distribution',
      'wiz.dir': 'Folder to clone your repos into',
      'wiz.dir.hint': 'Paste the full path: the commands below are built with it. It is created if it does not exist.',
      'wiz.dir.how.windows': 'To copy it: open the folder in File Explorer, click the address bar and copy (Ctrl+C). Right-click → “Copy as path” works too.',
      'wiz.dir.how.macos': 'To copy it: select the folder in Finder and press Option+Cmd+C (or Option-right-click → “Copy as Pathname”). You can use <code>~</code> for your home folder.',
      'wiz.dir.how.linux': 'To copy it: press Ctrl+L in your file manager and copy the path, or run <code>pwd</code> in a terminal inside the folder. You can use <code>~</code> for your home folder.',
      'wiz.options': 'Options',
      'wiz.forks': 'Include forks',
      'wiz.archived': 'Include archived repos',
      'wiz.warn.empty': 'Enter a folder to continue.',
      'wiz.warn.relative': 'Use a full path, for example {ex}.',
      's1.t': 'Install git and GitHub CLI',
      's1.win': 'Open <strong>PowerShell</strong> and run (skip this if you already have them):',
      's1.win.after': 'Close PowerShell and open it again so it picks up the new commands.',
      's1.mac': 'With <a href="https://brew.sh">Homebrew</a>, in the <strong>Terminal</strong> app:',
      's1.linux': "With your distribution's package manager:",
      's1.linux.after': 'If your distribution does not ship <code>gh</code> or ships an old version, follow the <a href="{url}">official GitHub CLI guide</a>.',
      's2.t': 'Link your GitHub account',
      's2.d': 'Sign in once on this machine:',
      's2.answers': 'GitHub CLI will ask a few questions. Answer like this:',
      's2.q1': 'Where do you use GitHub?',
      's2.q2': 'What is your preferred protocol for Git operations on this host?',
      's2.q3': 'Authenticate Git with your GitHub credentials?',
      's2.q4': 'How would you like to authenticate GitHub CLI?',
      's2.browser': 'Copy the 8-character code it shows, press Enter and paste it into the GitHub page that opens. Then authorize access.',
      's2.verify': 'Check that it is linked (it should show your username):',
      's2.note': 'GitHub CLI stores the session on your machine. gitnomad never reads or stores your token: it asks <code>gh</code> to authenticate when needed.',
      's3.t': 'Install gitnomad',
      's3.win': 'In PowerShell:',
      's3.unix': 'In the terminal:',
      's3.unix.after': 'If it then says <code>gitnomad</code> is not found, open a new terminal.',
      's3.manual': 'Prefer a manual download? It is on <a href="{url}">Releases</a>.',
      's4.t': 'Set your repos folder',
      's4.d': "Save the folder in this machine's configuration:",
      's4.opts': 'The options you picked:',
      's4.note': 'Each machine has its own configuration, so another device can use a different path.',
      's5.t': 'Check and clone',
      's5.doctor': 'Check that everything is ready (and fix what it can):',
      's5.list': 'See your repos and which ones are missing:',
      's5.dry': 'Dry run, nothing is downloaded:',
      's5.clone': 'Clone what is missing:',
      's5.after': 'On another machine, repeat these steps. When you create new repos, run <code>gitnomad clone</code> again.',
      'copy': 'Copy',
      'copied': 'Copied',
      'copy.toast': 'Command copied',
      'copy.fail': 'Could not copy; select it manually',
      'safety.title': 'Built not to break anything',
      'safety.1.t': 'Never deletes or overwrites.',
      'safety.1.d': 'If a folder with a repo\'s name exists but is not that repo, it is flagged as a conflict and left alone.',
      'safety.2.t': 'No tokens in sight.',
      'safety.2.d': "Authentication is handled by GitHub CLI using GitHub's official sign-in.",
      'safety.3.t': 'Never hangs.',
      'safety.3.d': 'It never prompts for passwords in the console. If a repo fails, it reports it and moves on.',
      'safety.4.t': 'Verified installer.',
      'safety.4.d': 'Every download is checked against its SHA-256 checksum before it is installed.',
      'cmds.title': 'All commands',
      'cmds.lead': 'Quick reference. <code>clone</code> options apply to that run only; to keep them, use <code>gitnomad config set</code>.',
      'cmds.keys': '<code>config set</code> keys: <code>base_dir</code>, <code>parallel</code>, <code>include_forks</code>, <code>include_archived</code> and <code>exclude</code>.',
      cmds: [
        ['gitnomad doctor', 'Checks that git, gh and your GitHub session are ready.'],
        ['gitnomad doctor --fix', 'Also fixes what it can, such as setting git to use the gh credentials.'],
        ['gitnomad config show', 'Shows your configuration and where the file is.'],
        ['gitnomad config init --dir PATH', 'Creates this machine’s configuration. With <code>--pick</code> you choose the folder in a window; <code>--force</code> replaces an existing one.'],
        ['gitnomad config set base_dir --pick', 'Opens your system’s folder picker and saves the folder you choose.'],
        ['gitnomad config set <key> <value>', 'Changes an option, for example <code>parallel 8</code> or <code>exclude "test-*"</code>.'],
        ['gitnomad list', 'Lists your repos and whether each is downloaded, missing, in conflict or skipped. Add <code>--json</code> for scripts.'],
        ['gitnomad clone --dry-run', 'Shows what would be cloned and where, without downloading anything.'],
        ['gitnomad clone', 'Clones the missing repos. Options: <code>--parallel N</code> and <code>--dir PATH</code>.'],
        ['gitnomad version', 'Shows the installed version.']
      ],
      'faq.title': 'Frequently asked questions',
      'footer.by': 'made by',
      'footer.code': 'Code on GitHub',
      faq: [
        ['Which repos does it download?', 'The public and private repos of your account, forks included. Archived repos are skipped unless you turn the option on. Organization repos are not included yet.'],
        ['Does it update repos I already have?', 'For now it clones what is missing and leaves existing repos alone. Syncing (pulling changes and flagging unpushed commits) is coming in the next version.'],
        ['Is the command-line output in English?', 'Not yet: gitnomad currently prints its messages in Spanish. English output is planned. The commands and options are the same.'],
        ['Does it work with several GitHub accounts?', 'It uses the active GitHub CLI account. Switch with <code>gh auth switch</code> and point to another folder with <code>gitnomad config set base_dir</code>.'],
        ['Where is the configuration stored?', '<p>Windows: <code>%APPDATA%\\gitnomad\\config.json</code></p><p>macOS: <code>~/Library/Application Support/gitnomad/config.json</code></p><p>Linux: <code>~/.config/gitnomad/config.json</code></p>'],
        ['Can I choose the folder in a window instead of typing it?', 'Yes: <code>gitnomad config set base_dir --pick</code> opens your system’s folder picker. On Windows and macOS it works out of the box; on Linux it needs <code>zenity</code> or <code>kdialog</code>. If you close the window without choosing, nothing changes. Available from v0.2.0 (check yours with <code>gitnomad version</code>).'],
        ['Can I clone into a folder that already has other projects?', 'Yes. Each repo goes into its own subfolder, repos already there are recognized as downloaded and nothing existing is modified. If a folder has a repo’s name but is something else, it is marked as a conflict and left alone. Run <code>gitnomad clone --dry-run</code> first to see what would happen.'],
        ['How do I exclude some repos?', 'With comma-separated patterns: <code>gitnomad config set exclude "test-*,old"</code>.'],
        ['How do I update gitnomad?', 'Run the install command again: it replaces the executable with the latest version.'],
        ['How do I uninstall it?', 'Delete the executable (Windows: <code>%LOCALAPPDATA%\\Programs\\gitnomad</code>; macOS/Linux: <code>~/.local/bin/gitnomad</code>) and the configuration folder. Your cloned repos are not touched.']
      ]
    }
  };

  // ---- Estado (se recuerda en este navegador si se puede) ----
  function load(key, fallback) {
    try { var v = localStorage.getItem('gitnomad.' + key); return v === null ? fallback : v; } catch (e) { return fallback; }
  }
  function save(key, value) {
    try { localStorage.setItem('gitnomad.' + key, value); } catch (e) { /* sin almacenamiento */ }
  }

  function detectOS() {
    var p = ((navigator.userAgentData && navigator.userAgentData.platform) || navigator.platform || navigator.userAgent || '').toLowerCase();
    if (p.indexOf('win') >= 0) return 'windows';
    if (p.indexOf('mac') >= 0 || p.indexOf('iphone') >= 0 || p.indexOf('ipad') >= 0) return 'macos';
    if (p.indexOf('linux') >= 0 || p.indexOf('x11') >= 0) return 'linux';
    return 'windows';
  }
  function detectLang() {
    var l = (navigator.language || 'es').toLowerCase();
    return l.indexOf('es') === 0 ? 'es' : 'en';
  }

  var state = {
    lang: load('lang', detectLang()),
    os: load('os', detectOS()),
    distro: load('distro', 'apt'),
    dir: load('dir', ''),
    forks: load('forks', '1') === '1',
    archived: load('archived', '0') === '1'
  };
  if (!T[state.lang]) state.lang = 'es';
  if (!DEFAULT_DIR[state.os]) state.os = 'windows';
  if (!state.dir) state.dir = DEFAULT_DIR[state.os];

  function t(key, vars) {
    var s = T[state.lang][key];
    if (s === undefined) s = T.es[key] || key;
    if (vars) Object.keys(vars).forEach(function (k) { s = s.split('{' + k + '}').join(vars[k]); });
    return s;
  }

  function esc(s) {
    return String(s).replace(/[&<>"']/g, function (c) {
      return { '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' }[c];
    });
  }

  // Comillas seguras para cada shell.
  function quotePS(s) { return "'" + s.replace(/'/g, "''") + "'"; }
  function quoteSh(s) { return "'" + s.replace(/'/g, "'\\''") + "'"; }

  function isAbsolute(dir, os) {
    if (dir === '~' || /^~[\\/]/.test(dir)) return true;
    if (os === 'windows') return /^[a-zA-Z]:[\\/]/.test(dir) || /^\\\\/.test(dir);
    return dir.charAt(0) === '/';
  }

  // ---- Bloques de HTML ----
  function cmd(text, prompt) {
    return '<div class="cmd"><pre data-prompt="' + esc(prompt) + '"><code>' + esc(text) +
      '</code></pre><button type="button" class="copy" data-copy="' + esc(text) + '">' + esc(t('copy')) + '</button></div>';
  }
  function p(html, cls) { return '<p' + (cls ? ' class="' + cls + '"' : '') + '>' + html + '</p>'; }
  function step(title, body) { return '<li class="step"><h3>' + esc(title) + '</h3>' + body + '</li>'; }

  function renderSteps() {
    var win = state.os === 'windows';
    var prompt = win ? 'PS>' : '$';
    var q = win ? quotePS : quoteSh;
    var dir = state.dir.trim();
    var html = '';

    // 1. Requisitos
    var s1 = '';
    if (win) {
      s1 += p(t('s1.win'));
      s1 += cmd('winget install --id Git.Git -e', prompt);
      s1 += cmd('winget install --id GitHub.cli -e', prompt);
      s1 += p(t('s1.win.after'), 'note');
    } else if (state.os === 'macos') {
      s1 += p(t('s1.mac'));
      s1 += cmd('brew install git gh', prompt);
    } else {
      s1 += p(t('s1.linux'));
      var linuxCmd = { apt: 'sudo apt update && sudo apt install -y git gh', dnf: 'sudo dnf install -y git gh', pacman: 'sudo pacman -S --needed git github-cli' }[state.distro];
      s1 += cmd(linuxCmd, prompt);
      s1 += p(t('s1.linux.after', { url: GH_INSTALL_URL }), 'note');
    }
    html += step(t('s1.t'), s1);

    // 2. Vincular GitHub
    var s2 = p(t('s2.d')) + cmd('gh auth login', prompt) + p(t('s2.answers'));
    s2 += '<ul class="answers">' +
      [['s2.q1', 'GitHub.com'], ['s2.q2', 'HTTPS'], ['s2.q3', 'Yes'], ['s2.q4', 'Login with a web browser']]
        .map(function (qa) { return '<li><span class="q">? ' + esc(t(qa[0])) + '</span> → <span class="a">' + esc(qa[1]) + '</span></li>'; }).join('') +
      '</ul>';
    s2 += p(t('s2.browser'));
    s2 += p(t('s2.verify')) + cmd('gh auth status', prompt);
    s2 += p(t('s2.note'), 'note');
    html += step(t('s2.t'), s2);

    // 3. Instalar gitnomad
    var s3 = '';
    if (win) {
      s3 += p(t('s3.win')) + cmd('irm ' + CONFIG.site + '/install.ps1 | iex', prompt);
    } else {
      s3 += p(t('s3.unix')) + cmd('curl -fsSL ' + CONFIG.site + '/install.sh | sh', prompt);
      s3 += p(t('s3.unix.after'), 'note');
    }
    s3 += p(t('s3.manual', { url: RELEASES_URL }));
    html += step(t('s3.t'), s3);

    // 4. Carpeta
    var s4 = p(t('s4.d'));
    s4 += cmd('gitnomad config set base_dir ' + q(dir || DEFAULT_DIR[state.os]), prompt);
    var extra = [];
    if (!state.forks) extra.push('gitnomad config set include_forks false');
    if (state.archived) extra.push('gitnomad config set include_archived true');
    if (extra.length) {
      s4 += p(t('s4.opts'));
      extra.forEach(function (c) { s4 += cmd(c, prompt); });
    }
    s4 += p(t('s4.note'), 'note');
    html += step(t('s4.t'), s4);

    // 5. Revisar y clonar
    var s5 = p(t('s5.doctor')) + cmd('gitnomad doctor --fix', prompt) +
      p(t('s5.list')) + cmd('gitnomad list', prompt) +
      p(t('s5.dry')) + cmd('gitnomad clone --dry-run', prompt) +
      p(t('s5.clone')) + cmd('gitnomad clone', prompt) +
      p(t('s5.after'), 'note');
    html += step(t('s5.t'), s5);

    document.getElementById('steps').innerHTML = html;
  }

  function validateDir() {
    var dir = state.dir.trim();
    var warn = document.getElementById('dir-warn');
    var msg = '';
    if (!dir) msg = t('wiz.warn.empty');
    else if (!isAbsolute(dir, state.os)) msg = t('wiz.warn.relative', { ex: '<code>' + esc(DEFAULT_DIR[state.os]) + '</code>' });
    warn.innerHTML = msg;
    warn.hidden = !msg;
  }

  function renderFAQ() {
    document.getElementById('faq-list').innerHTML = T[state.lang].faq.map(function (item) {
      var body = item[1].indexOf('<p>') === 0 ? item[1] : '<p>' + item[1] + '</p>';
      return '<details><summary>' + esc(item[0]) + '</summary><div>' + body + '</div></details>';
    }).join('');
  }

  function renderCommands() {
    document.getElementById('cmd-list').innerHTML = T[state.lang].cmds.map(function (c) {
      return '<li><code>' + esc(c[0]) + '</code><span>' + c[1] + '</span></li>';
    }).join('');
  }

  function applyLang() {
    document.documentElement.lang = state.lang;
    document.querySelectorAll('[data-i18n]').forEach(function (el) {
      el.innerHTML = t(el.getAttribute('data-i18n'));
    });
    var meta = document.querySelector('meta[name="description"]');
    if (meta) meta.setAttribute('content', t('meta.description'));
    document.querySelectorAll('[data-lang]').forEach(function (b) {
      b.setAttribute('aria-pressed', String(b.getAttribute('data-lang') === state.lang));
    });
    renderCommands();
    renderFAQ();
    renderSteps();
    validateDir();
  }

  function syncForm() {
    document.querySelectorAll('input[name="os"]').forEach(function (r) { r.checked = r.value === state.os; });
    document.getElementById('distro-field').hidden = state.os !== 'linux';
    document.getElementById('distro').value = state.distro;
    document.getElementById('dir').value = state.dir;
    var how = document.getElementById('dir-how');
    how.setAttribute('data-i18n', 'wiz.dir.how.' + state.os);
    how.innerHTML = t(how.getAttribute('data-i18n'));
    document.getElementById('dir').placeholder = DEFAULT_DIR[state.os];
    document.getElementById('forks').checked = state.forks;
    document.getElementById('archived').checked = state.archived;
  }

  // ---- Copiar ----
  var toastTimer;
  function toast(msg) {
    var el = document.getElementById('toast');
    el.textContent = msg;
    el.classList.add('show');
    clearTimeout(toastTimer);
    toastTimer = setTimeout(function () { el.classList.remove('show'); }, 1600);
  }
  function fallbackCopy(text) {
    var ta = document.createElement('textarea');
    ta.value = text;
    ta.setAttribute('readonly', '');
    ta.style.position = 'fixed';
    ta.style.opacity = '0';
    document.body.appendChild(ta);
    ta.select();
    var ok = false;
    try { ok = document.execCommand('copy'); } catch (e) { ok = false; }
    document.body.removeChild(ta);
    return ok;
  }
  function copy(btn) {
    var text = btn.getAttribute('data-copy');
    var done = function () {
      btn.textContent = t('copied');
      btn.classList.add('done');
      toast(t('copy.toast'));
      setTimeout(function () { btn.textContent = t('copy'); btn.classList.remove('done'); }, 1600);
    };
    var failed = function () { toast(t('copy.fail')); };
    if (navigator.clipboard && window.isSecureContext) {
      navigator.clipboard.writeText(text).then(done, function () { fallbackCopy(text) ? done() : failed(); });
    } else {
      fallbackCopy(text) ? done() : failed();
    }
  }

  // ---- Eventos ----
  document.addEventListener('click', function (e) {
    var btn = e.target.closest('.copy');
    if (btn) { copy(btn); return; }
    var lb = e.target.closest('[data-lang]');
    if (lb) { state.lang = lb.getAttribute('data-lang'); save('lang', state.lang); applyLang(); }
  });

  document.querySelectorAll('input[name="os"]').forEach(function (r) {
    r.addEventListener('change', function () {
      var wasDefault = state.dir.trim() === '' || state.dir === DEFAULT_DIR[state.os];
      state.os = r.value;
      if (wasDefault) state.dir = DEFAULT_DIR[state.os];
      save('os', state.os); save('dir', state.dir);
      syncForm(); renderSteps(); validateDir();
    });
  });
  document.getElementById('distro').addEventListener('change', function (e) {
    state.distro = e.target.value; save('distro', state.distro); renderSteps();
  });
  document.getElementById('dir').addEventListener('input', function (e) {
    // «Copiar como ruta de acceso» de Windows pega la ruta entre comillas.
    var v = e.target.value.trim().replace(/^(["'])(.*)\1$/, '$2');
    if (v !== e.target.value.trim()) e.target.value = v;
    state.dir = v; save('dir', state.dir); renderSteps(); validateDir();
  });
  document.getElementById('forks').addEventListener('change', function (e) {
    state.forks = e.target.checked; save('forks', state.forks ? '1' : '0'); renderSteps();
  });
  document.getElementById('archived').addEventListener('change', function (e) {
    state.archived = e.target.checked; save('archived', state.archived ? '1' : '0'); renderSteps();
  });

  document.querySelectorAll('.js-repo-link').forEach(function (a) { a.href = REPO_URL; });

  syncForm();
  applyLang();
})();
