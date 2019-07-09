


numNodes = 1000

t = [0:1:1000]';

S0 = rand(1,numNodes)*0.2+0.1;
S1 = rand(1,numNodes)*0.2+0.1;
S2 = rand(1,numNodes)*0.2+0.1;

E0 = rand(1,numNodes)*0.1.*S0;
E1 = rand(1,numNodes)*0.1.*S1;
E2 = rand(1,numNodes)*0.1.*S2;


Tau1 = 10;
Tau2 = 500;
ET1 = Tau1*rand(1,numNodes)*0.05;
ET2 = Tau1*rand(1,numNodes)*0.05;


S = (S0+E0)+(S1+E1).*exp(-t./(Tau1+ET1))+(S2+E2).*exp(-t./(Tau2+ET2));

plot(S)