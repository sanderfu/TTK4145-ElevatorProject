# Comments

## 1
### main.go
* Ordermanager funksjonen burde ikke ta inn argumenter -> fjerner nytten av en 
config fil
* Samme som over med watchdog sendernode
* Watchdog.SenderNode burde endre navn

## 2
* Vanskelig å se hvilke komponeneter som er koblet til hvilke moduler
* Burde vi gjøre våre channels til å ta type `interface {}`, og så har vi bare 
slik

## 3 Functionality
* oppdater "read me" filene så de er ryddig og forteller kort hva "packagene" gjør
og hvordan systemet er satt opp (at det er mesh network)
* Alle "packagene" må ha egen "read me" file

## 4 Coherence
* queue.go - bør vi ha funksjonene OrderToTakeAtFloor(...) og backupListener i en annen fil, hvis de skal bli bør de endre navn?
* Noen andre bør også dobbeltsjekke dette :))

## 5 Completeness
* 