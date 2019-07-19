%generated = []
%comsol = []

x = 0:0.05:.95;

hold on
plot(x*3,comsol')

yyaxis right
plot(x*3,generated')

legend("Comsol Data", "Generated Data")

xlabel("Distance (m)")
ylabel("Raw Concentration")
