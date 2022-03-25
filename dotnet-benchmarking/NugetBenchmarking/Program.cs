using System;
using Serilog;

namespace NugetBenchmarking
{
    class Program
    {
        static void Main(string[] args)
        {
            // Serilog
            using var log = new LoggerConfiguration()
                  .WriteTo.Console()
                  .CreateLogger();
            log.Information("Hello, World!");
        }
    }
}
