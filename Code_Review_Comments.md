# Comments

## 1
### main.go
* Ordermanager funksjonen burde ikke ta inn argumenter -> fjerner nytten av en 
config fil
* Samme som over med watchdog ElevatorNode
* Watchdog.ElevatorNode burde endre navn

## 2
* Vanskelig å se hvilke komponeneter som er koblet til hvilke moduler
* Burde vi gjøre våre channels til å ta type `interface {}`, og så har vi bare 
slik
    * Nei, systemet er ikke designed for å kunne håndtere det, ville krevd en modul som bestemte men som skulle få meldingen eller lignende.

## 3 Functionality
* oppdater "read me" filene så de er ryddig og forteller kort hva "packagene" gjør
og hvordan systemet er satt opp (at det er mesh network)
* Alle "packagene" må ha egen "read me" file

## 4 Coherence
* queue.go - bør vi ha funksjonene OrderToTakeAtFloor(...) og backupTimeoutListener() 
i en annen fil, hvis de skal bli bør de endre navn?
* Noen andre bør også dobbeltsjekke dette :))

## 5 Completeness
* Noen andre bør også dobbelt sjekke om vi mangler noen funksjonalitet i modulene

## 6 State 
* Kanskje samle mye av det som skjer i moving-state i en ny state `servingOrder`
eller noe sånt
* Test om dette fungerer før det gjøres! Usikker på hva forskjellen på de to som 
sier at ordren var ferdig utgjør

## 7 Functions
* Skriv 80 //////// rundt de private globale variablene vi bruker
* Dobbeltsjekk at ingen funksjoner gjør to ting

## 8 Understandability
* Sjekk om funksjonene er enkle å forstå -> legg inn små kommentarer der det 
trengs
* Sjekk at alle variabler har gode og tydelige navn
* Sjekk på at vi ikke har for mange nestede løkker

## 9 Traceability
* Funksjoner tar inn channels som argument. Disse channelene burde ha navn som
viser hvilken modul som mottar de
* Alle kanaler har navn som inneholder BÅDE hvor de kommer fra OG hvor de skal til

## 10 Direction
* Vet ikke helt hva vi skal endre på for å gjøre dette bedre...

## 11 Comments
* Når man går gjennom funksjoner, legg inn kommentarer der trengs
* Gå over og sjekk at det ikke finnes unødvendige kommentarer

## 12 Naming
* Sjekk at funksjoner har bra navn som i hvert fall gir en indikasjon på hvilken
modul de hører til

## Generelt:
* Legg til 80 /////////////// der det trengs

## Oppdeling

* cmd
    * Elevator main (vanlig) (Martin)
    * Watchdog main (Tobias)
    * startElevator main (Sander)

* internal
    * Channels (felles)
    * configuration (Martin) (inkl. config.json)
    * datatypes (Tobias)
    * FSM (Sander)
    * HWmanager (Martin)
    * Networkmanager (Sander)
    * Ordermanager
        * Order (Sander)
        * Queue (Martin)
        * Cost (Tobias)
    * Wathcdog (Tobias)