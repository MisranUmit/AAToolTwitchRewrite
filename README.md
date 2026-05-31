# AAToolTwitchRewrite - Minecraft Advancement Monitor

WIP advancement monitor for AA speedruns.

## Systemkrav

- **Go 1.26.1 eller nyere**
- **Minecraft-installasjon**

### Bygging fra kildekode

```bash
# Klon eller last ned prosjektet
cd AAToolTwitchRewrite

# Last ned avhengigheter
go mod download

# Kompiler programmet
go build -o AAToolMonitor.exe

# Eller bare kjør direkte
go run main.go advancementData.go
```

### Output

Resultatet blir `AAToolMonitor.exe` i samme mappe.

## ⚙️ Konfigurasjon

`config.json`-fil i samme mappe som programmet.

### Første gangs oppstart

Når du kjører programmet for første gang:

1. Opprett eller rediger `config.json`:
```json
{
  "minecraftPath": "C:\\<INSTANCE_MAPPE>\\minecraft
}
```

2. Stien skal peke til `minecraft`-mappen
   - **Offisiell launcher**: `C:\Users\[Brukernavn]\AppData\Roaming\.minecraft`
   - **Prism Launcher**: `C:\Users\[Brukernavn]\AppData\Roaming\PrismLauncher\instances\[Instansnavn]\minecraft`

## Bruk

```bash
# Kjør programmet
./AAToolMonitor.exe

# Eller fra PowerShell
.\AAToolMonitor.exe

# Eller direkte med Go
go run main.go advancementData.go
```

### Hva som skjer

1. Programmet starter og validerer Minecraft-stien fra `config.json`
2. Den finner den mest nylig modifiserte verden (prefererer verdener med advancements-data)
3. Laster inn alle eksisterende advancements
4. Venter på nye advancements mens du spiller
5. Når du oppnår et nytt advancement, vises det med navn som: `NEW ADVANCEMENT: Stone Age`

### Eksempel på log

```
Monitoring world: World_1
Watching: C:\Users\Scott\AppData\Roaming\.minecraft\saves\World_1\advancements
Loaded 15 initial advancements
Waiting for new advancements...
NEW ADVANCEMENT: Stone Age
NEW ADVANCEMENT: Getting an Upgrade
NEW ADVANCEMENT: Acquire Hardware
```