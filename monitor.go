package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/mem"
)

func main() {
	for {
		// Obtém a temperatura da CPU
		temp, err := getCPUTemperature()
		if err != nil {
			fmt.Println("Erro ao obter temperatura da CPU:", err)
		}

		// Obtém o uso de memória
		memUsage, err := getMemoryUsage()
		if err != nil {
			fmt.Println("Erro ao obter uso de memória:", err)
		}

		// Formata a mensagem
		message := fmt.Sprintf("CPU Temp: %.2f°C\nMem Usage: %.2f%%", temp, memUsage)

		// Exibe a notificação
		err = showNotification(message)
		if err != nil {
			fmt.Println("Erro ao exibir notificação:", err)
		}

		// Aguarda 5 segundos antes de repetir
		time.Sleep(5 * time.Second)
	}
}

func getCPUTemperature() (float64, error) {
	// Executa o comando para obter a temperatura
	cmd := exec.Command("sensors")
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	// Procura pela linha com a temperatura
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Core") || strings.Contains(line, "Package id 0:") {
			fields := strings.Fields(line)
			for _, field := range fields {
				if strings.Contains(field, "+") && strings.Contains(field, "°C") {
					tempStr := strings.TrimSuffix(strings.TrimPrefix(field, "+"), "°C")
					temp, err := strconv.ParseFloat(tempStr, 64)
					if err != nil {
						return 0, err
					}
					return temp, nil
				}
			}
		}
	}
	return 0, fmt.Errorf("Temperatura da CPU não encontrada")
}

func getMemoryUsage() (float64, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return 0, err
	}
	return v.UsedPercent, nil
}

func showNotification(message string) error {
	cmd := exec.Command("notify-send", "Monitoramento do Sistema", message)
	return cmd.Run()
}
