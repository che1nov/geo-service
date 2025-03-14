FROM grafana/grafana:latest

# Копируем datasource и дашборд
COPY datasources.yml /etc/grafana/provisioning/datasources/
COPY dashboards/geo-dashboard.json /var/lib/grafana/dashboards/

# Настраиваем права
RUN chmod 644 /etc/grafana/provisioning/datasources/datasources.yml
RUN chmod 644 /var/lib/grafana/dashboards/geo-dashboard.json

# Устанавливаем плагины и пароль
ENV GF_INSTALL_PLUGINS=grafana-piechart-panel
ENV GF_SECURITY_ADMIN_PASSWORD=admin