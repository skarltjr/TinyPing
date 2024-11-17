package main

import (
	"github.com/sirupsen/logrus"
	"html/template"
	"int-status/internal"
	"int-status/internal/cache"
	"int-status/internal/config"
	"int-status/internal/manager"
	"int-status/internal/storage"
	"net/http"
	"os"
	"strings"
	"time"
	_ "time/tzdata"
)

const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>TINY PING</title>
    <style>
        body {
            background-color: #1a1a1a;
            color: white;
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Arial, sans-serif;
            margin: 0;
            padding: 20px;
        }
        h1 {
            text-align: center;
            font-size: 2.5em;
            margin-bottom: 40px;
            font-weight: normal;
        }

        /* Incidents 섹션 스타일 */
        .incidents-section {
            max-width: 1200px;
            margin: 0 auto 40px;
            padding: 0 20px;
            box-sizing: border-box;
        }
        .incidents-title {
            font-size: 1.5em;
            margin-bottom: 20px;
            font-weight: normal;
        }
        .incident-card {
            background-color: rgba(255, 255, 255, 0.05);
            border-radius: 12px;
            padding: 20px;
            margin-bottom: 12px;
        }
        .incident-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 4px;
        }
        .incident-service {
            font-size: 1.2em;
        }
        .incident-time {
            color: rgba(255, 255, 255, 0.5);
            font-size: 0.9em;
        }
        .no-incidents {
            text-align: center;
            padding: 30px;
            background-color: rgba(76, 175, 80, 0.05);
            border-radius: 12px;
            border: 1px solid rgba(76, 175, 80, 0.1);
        }
        .perfect-day {
            color: #4CAF50;
            font-size: 1.5em;
            font-weight: 500;
            margin-bottom: 8px;
        }
        .sub-message {
            color: rgba(255, 255, 255, 0.7);
            font-size: 0.95em;
        }
        .section-divider {
            border: none;
            border-top: 1px solid rgba(255, 255, 255, 0.1);
            margin: 40px auto;
            max-width: 1200px;
        }

        /* Services 섹션 스타일 */
        .dashboard {
            display: grid;
            grid-template-columns: repeat(3, minmax(0, 1fr));
            gap: 20px;
            width: 100%;
            max-width: 1200px;
            margin: 0 auto;
            padding: 0 20px;
            box-sizing: border-box;
        }
        .service-card {
            background-color: rgba(255, 255, 255, 0.05);
            border-radius: 12px;
            padding: 24px;
            display: flex;
            flex-direction: column;
            gap: 8px;
            width: 100%;
            box-sizing: border-box;
        }
        .service-name {
            font-size: 1.8em;
            margin-bottom: 10px;
            font-weight: normal;
        }
        .service-status {
            display: flex;
            align-items: center;
            justify-content: space-between;
            gap: 10px;
            width: 100%;
        }
        .status-info {
            display: flex;
            flex-direction: column;
        }
        .status-text {
            font-size: 1.2em;
            margin-bottom: 4px;
        }
        .latency-text {
            color: rgba(255, 255, 255, 0.5);
            font-size: 0.9em;
        }
        .status-dots {
            display: flex;
            gap: 6px;
            align-items: center;
        }
        .dot {
            width: 10px;
            height: 10px;
            border-radius: 50%;
            position: relative;
            cursor: pointer;
        }
        .dot:hover::after {
            content: attr(data-timestamp);
            position: absolute;
            bottom: 100%;
            left: 50%;
            transform: translateX(-50%);
            background-color: rgba(0, 0, 0, 0.8);
            color: white;
            padding: 4px 8px;
            border-radius: 4px;
            font-size: 0.8em;
            white-space: nowrap;
            margin-bottom: 8px;
            z-index: 1;
        }
        .dot-up {
            background-color: #4CAF50;
        }
        .dot-down {
            background-color: #f44336;
        }
        .status-up {
            color: #4CAF50;
        }
        .status-down {
            color: #f44336;
        }

        @media (max-width: 1024px) {
            .dashboard {
                grid-template-columns: repeat(2, 1fr);
            }
        }
        @media (max-width: 768px) {
            .dashboard {
                grid-template-columns: 1fr;
            }
        }
    </style>
</head>
<body>
	<h1>
		<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 300 60">
			<defs>
				<linearGradient id="titleGradient" x1="0%" y1="0%" x2="100%" y2="0%">
					<stop offset="0%" style="stop-color:#4CAF50;"/>
					<stop offset="100%" style="stop-color:#2196F3;"/>
				</linearGradient>
				<filter id="titleShadow">
					<feDropShadow dx="1" dy="1" stdDeviation="1" flood-color="#000"/>
				</filter>
			</defs>
			
			<text x="50%" y="50%" text-anchor="middle" dominant-baseline="middle"
				  font-size="16" font-weight="bold" fill="url(#titleGradient)"
				  filter="url(#titleShadow)">
				TINY PING
			</text>
		</svg>
	</h1>
    <!-- Incidents 섹션 -->
    <div class="incidents-section">
        <h2 class="incidents-title">Today's Outages</h2>
        {{if .Incidents}}
            {{range $service, $serviceIncidents := .Incidents}}
                {{range $incident := $serviceIncidents}}
                <div class="incident-card">
                    <div class="incident-header">
                        <div class="incident-service">{{$service}}</div>
                        <div class="incident-time">Down: {{formatTime $incident.StartTime}} - {{formatTime $incident.EndTime}}</div>
                    </div>
                </div>
                {{end}}
            {{end}}
        {{else}}
            <div class="no-incidents">
                <div class="perfect-day">Perfect Day!</div>
                <div class="sub-message">All systems have been operational today</div>
            </div>
        {{end}}
    </div>

    <hr class="section-divider">

    <!-- Services 섹션 -->
    <div class="dashboard">
        {{range $service, $statuses := .Services}}
        <div class="service-card">
            <div class="service-name">{{$service}}</div>
            <div class="service-status">
                <div class="status-info">
                    <div class="status-text {{if eq (index $statuses (sub (len $statuses) 1)).Status "UP"}}status-up{{else}}status-down{{end}}">
                        {{if eq (index $statuses (sub (len $statuses) 1)).Status "UP"}}Operational{{else}}Down{{end}}
                    </div>
                    <div class="latency-text">{{(index $statuses (sub (len $statuses) 1)).Latency}} ms</div>
                </div>
                <div class="status-dots">
                    {{range $statuses}}
                    <div class="dot {{if eq .Status "UP"}}dot-up{{else}}dot-down{{end}}" 
                         data-timestamp="{{formatTime .Timestamp}}">
                    </div>
                    {{end}}
                </div>
            </div>
        </div>
        {{end}}
    </div>
</body>
</html>
`

const errorTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>TINY PING - Loading</title>
    <style>
        body {
            background-color: #1a1a1a;
            color: white;
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Arial, sans-serif;
            margin: 0;
            padding: 20px;
            display: flex;
            align-items: center;
            justify-content: center;
            min-height: 100vh;
        }
        .error-container {
            text-align: center;
            background-color: rgba(255, 255, 255, 0.05);
            border-radius: 12px;
            padding: 40px;
            max-width: 500px;
            width: 100%;
        }
        .loading-spinner {
            border: 3px solid rgba(255, 255, 255, 0.1);
            border-radius: 50%;
            border-top-color: #2196F3;
            width: 40px;
            height: 40px;
            animation: spin 1s linear infinite;
            margin: 0 auto 20px;
        }
        @keyframes spin {
            to { transform: rotate(360deg); }
        }
        h1 {
            font-size: 1.8em;
            margin-bottom: 20px;
            font-weight: normal;
        }
        .message {
            color: #888;
            line-height: 1.6;
            margin-bottom: 20px;
        }
        .auto-refresh {
            color: #666;
            font-size: 0.9em;
        }
    </style>
    <script>
        setTimeout(function() {
            window.location.reload();
        }, 5000);
    </script>
</head>
<body>
    <div class="error-container">
        <div class="loading-spinner"></div>
        <h1>Loading Dashboard</h1>
        <div class="message">
            Please wait a moment while we gather the service information.
            The page will refresh automatically.
        </div>
        <div class="auto-refresh">
            Refreshing in 5 seconds...
        </div>
    </div>
</body>
</html>
`

var timeZoneLoc *time.Location

var funcMap = template.FuncMap{
	"sub": func(a, b int) int {
		return a - b
	},
	"formatTime": func(t time.Time) string {
		kstTime := t.In(timeZoneLoc)
		return kstTime.Format("2006-01-02 15:04:05")
	},
}

type DashboardData struct {
	Services  map[string][]internal.Status
	Incidents map[string][]internal.Incident
}

func main() {
	location, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		logrus.Fatalf("Error loading KST timezone: %v", err)
	}
	timeZoneLoc = location

	const yamlPath = "./config/config.yaml"
	serviceConfigs, err := config.LoadServices(yamlPath)
	if err != nil {
		logrus.Fatalf("Error loading services: %v", err)
	}

	region := GetEnv("AWS_REGION")
	tableName := GetEnv("DYNAMODB_TABLE_NAME")
	dbStorage, err := storage.NewDynamoDBStorage(region, tableName)
	if err != nil {
		logrus.Fatal(err)
	}

	serviceManager := manager.NewServiceManager(serviceConfigs, dbStorage)
	htmlCache := cache.NewHTMLCache(10 * time.Second)

	go func() {
		dashboardTmpl := template.Must(template.New("dashboard").Funcs(funcMap).Parse(htmlTemplate))
		errorTmpl := template.Must(template.New("error").Parse(errorTemplate))

		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if content, ok := htmlCache.Get(); ok {
				w.Header().Set("Content-Type", "text/html")
				w.Write([]byte(content))
				return
			}

			statuses, err := serviceManager.GetDailyServiceStatus()
			if err != nil {
				logrus.Errorf("Error getting service statuses: %v", err)
				w.Header().Set("Content-Type", "text/html")
				errorTmpl.Execute(w, nil)
				return
			}

			incidents, err := serviceManager.GetDailyIncidents()
			if err != nil {
				logrus.Errorf("Error getting service incidents: %v", err)
			}

			data := DashboardData{
				Services:  statuses,
				Incidents: incidents,
			}

			var buf strings.Builder
			if err := dashboardTmpl.Execute(&buf, data); err != nil {
				logrus.Errorf("Error executing template: %v", err)
				errorTmpl.Execute(w, nil)
				return
			}

			rendered := buf.String()
			htmlCache.Set(rendered)

			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(rendered))
		})

		logrus.Info("Starting server on :8080")
		logrus.Fatal(http.ListenAndServe(":8080", nil))
	}()

	serviceManager.StartMonitoring(1 * time.Minute)
}

func GetEnv(key string) string {
	return os.Getenv(key)
}
