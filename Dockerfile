FROM mcr.microsoft.com/dotnet/sdk:7.0-alpine as builder

WORKDIR /app

COPY ReminderBot.Bot/ReminderBot.csproj ReminderBot.csproj

RUN dotnet restore

COPY ReminderBot.Bot .

RUN dotnet publish -c Release

FROM mcr.microsoft.com/dotnet/runtime:7.0-alpine as runner

COPY --from=builder /app/bin/Release/net7.0/publish /publish

ENTRYPOINT [ "dotnet", "/publish/ReminderBot.dll" ]
