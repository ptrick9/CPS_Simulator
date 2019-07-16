


numNodes = 1000

t = [0:1:10000]';

S0 = rand(1,numNodes)*0.2+0.1;
S1 = rand(1,numNodes)*0.2+0.1;
S2 = rand(1,numNodes)*0.2+0.1;

E0 = rand(1,numNodes)*0.1.*S0;
E1 = rand(1,numNodes)*0.1.*S1;
E2 = rand(1,numNodes)*0.1.*S2;

%E0 = 0;
%E1 = 0;
%E2 = 0;

Tau1 = 1000;
Tau2 = 1000;
ET1 = Tau1*rand(1,numNodes)*0.05;
ET2 = Tau1*rand(1,numNodes)*0.05;


S = (S0+E0)+(S1+E1).*exp(-t./(Tau1+ET1))+(S2+E2).*exp(-t./(Tau2+ET2));

plot(S)


half = zeros(1,numNodes);
long = 0;
for i=1:1:numNodes
    node = S(:,i)';
    a = node <= node(1)/2;
    ind = find(a, 1, 'first');
    if isempty(ind) 
        half(i) = 0;
        long = long + 1;
    else
        half(i) = ind;
    end

end

mean(half)
long