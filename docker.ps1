
Do
{
    $prompt = @(
    "-----------------------------------------:"
    "Select operation:"
    '1      - docker-build'
    '2      - docker-up'
    '3      - docker-down'
    '4      - docker-restart(down-up)'
    '5      - docker-restart'

    'logs   - logs'
    'ps     - lists containers'
    'pss    - lists only services'
    'images - list images'
    'config - view docker-compose.yaml'
    'help   - view help'

    'pg     - postgres rebuild'
    'php    - php rebuild'
    'm      - mysql rebuild'
    'b      - rebuild all'
    "0      - exit "
    ) -join "`n "
    $operation = Read-Host $prompt

    if ($operation -eq '1')
    {
        docker compose build
    }
    if ($operation -eq '2')
    {
        docker compose up -d
    }
    if ($operation -eq '3')
    {
        docker compose down
    }
    if ($operation -eq '4')
    {
        docker compose down
        docker compose up -d
    }
    if ($operation -eq '5')
    {
        docker compose restart
    }

    if ($operation -eq 'logs')
    {
        docker compose logs
    }
    if ($operation -eq 'ps')
    {
        docker-compose ps
    }
    if ($operation -eq 'pss')
    {
        docker-compose ps --services
    }
    if ($operation -eq 'images')
    {
        docker-compose images
    }
    if ($operation -eq 'config')
    {
        docker-compose config
    }
    if ($operation -eq 'help')
    {
        docker-compose help
    }

    if ($operation -eq 'all')
    {
        docker compose down
        docker compose build
        docker compose up -d
    }
    if ($operation -eq 'pg')
    {
        docker compose down
        docker compose build postgres
        docker compose up -d
    }
    if ($operation -eq 'php')
    {
        docker compose down
        docker compose build php-8.0
        docker compose up -d
    }
    if ($operation -eq 'mysql')
    {
        docker compose down
        docker compose build mysql-8
        docker compose up -d
    }

    if ($operation -eq '0')
    {
        break
    }
}
While (1)


