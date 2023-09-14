FROM mcr.microsoft.com/dotnet/sdk:7.0-alpine as builder

ARG VERSION=0.1.0

WORKDIR /app

COPY ReminderBot.Bot/ReminderBot.csproj ReminderBot.csproj

RUN dotnet restore

COPY ReminderBot.Bot .

# Inject semver into app
RUN echo "namespace ReminderBot.Bot;" > BotVersion.cs
RUN echo "" >> BotVersion.cs
RUN echo "public class BotVersion {" >> BotVersion.cs
RUN echo "    public const string Version = \"${VERSION}\";" >> BotVersion.cs
RUN echo "}" >> BotVersion.cs
RUN echo "" >> BotVersion.cs

RUN dotnet publish -c Release /p:Version=${VERSION}

FROM mcr.microsoft.com/dotnet/runtime:7.0-alpine as runner

COPY --from=builder /app/bin/Release/net7.0/publish /publish

ENTRYPOINT [ "dotnet", "/publish/ReminderBot.dll" ]
