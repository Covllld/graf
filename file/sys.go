package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/StackExchange/wmi"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

type Win32_VideoController struct {
	Name               string
	AdapterRAM         uint64
	DriverVersion      string
	VideoProcessor     string
	CurrentRefreshRate uint16
	MaxRefreshRate     uint16
}

func getGPUInfo() (float64, float64, error) {
	var gpuInfo []Win32_VideoController
	query := "SELECT Name, AdapterRAM, DriverVersion, VideoProcessor FROM Win32_VideoController"
	err := wmi.Query(query, &gpuInfo)
	if err != nil {
		return 0, 0, err
	}

	// Essayer d'obtenir l'utilisation GPU via une commande système
	cmd := exec.Command("powershell", "-Command",
		"(Get-WmiObject Win32_PerfFormattedData_GPUPerformanceCounters_GPU).PercentProcessorTime")
	gpuUtilOutput, err := cmd.Output()
	gpuUtil := 0.0
	if err == nil {
		gpuUtil, _ = strconv.ParseFloat(strings.TrimSpace(string(gpuUtilOutput)), 64)
	}

	// La température GPU est difficile à obtenir universellement
	gpuTemp := 0.0
	tempCmd := exec.Command("powershell", "-Command",
		"(Get-WmiObject MSAcpi_ThermalZoneTemperature -Namespace root/wmi).CurrentTemperature / 10 - 273")
	gpuTempOutput, err := tempCmd.Output()
	if err == nil {
		gpuTemp, _ = strconv.ParseFloat(strings.TrimSpace(string(gpuTempOutput)), 64)
	}

	return gpuUtil, gpuTemp, nil
}

func getCPUTemperature() float64 {
	// Méthode alternative pour obtenir la température CPU
	cmd := exec.Command("powershell", "-Command",
		"(Get-WmiObject -Namespace root\\wmi -Class MSAcpi_ThermalZoneTemperature).CurrentTemperature / 10 - 273")
	output, err := cmd.Output()
	if err != nil {
		return 0
	}

	temp, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
	if err != nil {
		return 0
	}

	return temp
}

func main() {
	for {
		// Informations CPU
		cpuPercent, err := cpu.Percent(time.Second, false)
		if err != nil {
			log.Printf("Erreur CPU : %v", err)
		}

		// Obtenir la température du CPU via une méthode alternative
		cpuTemp := getCPUTemperature()

		// Informations GPU
		gpuUtil, gpuTemp, err := getGPUInfo()
		if err != nil {
			log.Printf("Erreur GPU : %v", err)
		}

		// Informations RAM
		vmStat, err := mem.VirtualMemory()
		if err != nil {
			log.Printf("Erreur mémoire : %v", err)
		}

		// Effacer le terminal (compatible multi-plateforme)
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()

		// Affichage des informations
		fmt.Println("===== Moniteur de Ressources Système =====")

		if cpuTemp > 0 {
			fmt.Printf("🌡️  Température CPU: %.1f°C\n", cpuTemp)
		}

		if len(cpuPercent) > 0 {
			fmt.Printf("💻 Pourcentage CPU: %.1f%%\n", cpuPercent[0])
		}

		if vmStat != nil {
			fmt.Printf("🧠 Utilisation RAM: %.1f%%\n", vmStat.UsedPercent)
		}

		fmt.Printf("🖥️  Température GPU: %.1f°C\n", gpuTemp)
		fmt.Printf("🔋 Pourcentage GPU: %.1f%%\n", gpuUtil)

		// Attendre avant le prochain rafraîchissement
		time.Sleep(2 * time.Second)
	}
}
