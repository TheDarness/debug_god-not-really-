package main

import (
	"fmt"
)

type Channel struct {
	name        string
	capacity    int
	currentLoad int
	overloaded  bool
}

type RoutingStats struct {
	totalTraffic       int
	routedTraffic      int
	droppedTraffic     int
	overloadedChannels int
	usedChannels       int
}

// createChannels створює список доступних каналів зв'язку
func createChannels() []Channel {
	channels := []Channel{
		{"Main", 1000, 0, false},
		{"Backup1", 600, 0, false},
		{"Backup2", 400, 0, false},
	}

	return channels
}

// createTrafficProfile повертає профіль трафіку по хвилинах
func createTrafficProfile() []int {
	return []int{500, 800, 300, 1200, 900, 400}
}

// routeMinute має розподілити трафік цієї хвилини по каналах
// і повернути: routed, dropped
func routeMinute(traffic int, channels []Channel) (int, int) {
	remaining := traffic
	routed := 0
	dropped := 0

	// поки є трафік, намагаємося кудись його запхати
	for remaining > 0 {
		// шукаємо "найменш завантажений" канал
		minIndex := -1
		maxLoad := -1

		for i := 0; i < len(channels); i++ {
			// пропускаємо вже перевантажені канали
			if channels[i].currentLoad >= channels[i].capacity {
				continue
			}
			if channels[i].currentLoad > maxLoad {
				maxLoad = channels[i].currentLoad
				minIndex = i
			}
		}

		// якщо не знайшли підходящий канал — викидаємо трафік
		if minIndex == -1 {
			dropped += remaining
			remaining = 0
			break
		}

		// обраний канал
		ch := &channels[minIndex]
		free := ch.capacity - ch.currentLoad
		take := remaining
		if take > free {
			take = free
		}
		ch.currentLoad += take
		remaining -= take
		routed += take

		// перевірка перевантаження
		if ch.currentLoad > ch.capacity {
			ch.overloaded = true
		}
	}

	return routed, dropped
}

// simulateRouting проходить по всіх хвилинах трафіку
func simulateRouting(traffic []int, channels []Channel, maxCapacity int) RoutingStats {
	stats := RoutingStats{}

	fmt.Println("Starting traffic routing simulation...")
	fmt.Println("Minutes:", len(traffic), "Channels:", len(channels))

	for i := 0; i < len(traffic); i++ {
		minuteTraffic := traffic[i]
		fmt.Println("\nMinute", i+1, "incoming traffic:", minuteTraffic)

		routed, dropped := routeMinute(minuteTraffic, channels)

		fmt.Println("Routed:", routed, "Dropped:", dropped)

		stats.totalTraffic += minuteTraffic
		stats.routedTraffic += routed
		stats.droppedTraffic += dropped
	}
	// рахуємо використані та перевантажені канали
	for i := 0; i < len(channels); i++ {
		if channels[i].currentLoad > 0 {
			stats.usedChannels++
		}
		if channels[i].overloaded {
			stats.overloadedChannels++
		}
	}

	return stats
}

func printReport(stats RoutingStats, channels []Channel) {
	fmt.Println("\n=== Traffic routing report ===")
	fmt.Println("Total traffic:", stats.totalTraffic)
	fmt.Println("Routed traffic:", stats.routedTraffic)
	fmt.Println("Dropped traffic:", stats.droppedTraffic)
	fmt.Println("Used channels:", stats.usedChannels)
	fmt.Println("Overloaded channels:", stats.overloadedChannels)

	fmt.Println("\nChannel states:")
	for _, ch := range channels {
		fmt.Println("-",
			ch.name,
			"load:", ch.currentLoad,
			"capacity:", ch.capacity,
			"overloaded:", ch.overloaded,
		)
	}

	if stats.droppedTraffic == 0 && stats.overloadedChannels > 0 {
		fmt.Println("Conclusion: Увесь трафік успішно розподілено, перевантажень не було.")
	} else if stats.droppedTraffic > 0 && stats.overloadedChannels == 0 {
		fmt.Println("Conclusion: Балансувальник захистив канали, але частину трафіку скинуто.")
	} else {
		fmt.Println("Conclusion: Виявлено перевантажені канали — алгоритм розподілу трафіку потребує перегляду.")
	}
}

func main() {
	channels := createChannels()
	traffic := createTrafficProfile()
	maxCapacity := 1000

	fmt.Println("Simulating network traffic distribution...")

	stats := simulateRouting(traffic, channels, maxCapacity)

	printReport(stats, channels)
}
